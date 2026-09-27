package main

import (
	"testing"

	"gin-blog/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSeedMenusSurviveTranslationAndRepairOrphans(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:seed_menu_repair?mode=memory&cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&model.Menu{}, &model.Role{}, &model.RoleMenu{}))

	role := model.Role{Name: "admin", Label: "Administrator"}
	assert.NoError(t, db.Create(&role).Error)

	oldParent := model.Menu{Name: "Permission", Path: "/auth", Component: "Layout"}
	duplicateParent := model.Menu{Name: "Permission Management", Path: "/auth", Component: "Layout"}
	assert.NoError(t, db.Create(&oldParent).Error)
	assert.NoError(t, db.Create(&duplicateParent).Error)

	oldChild := model.Menu{Name: "Menu", Path: "menu", Component: "/auth/menu", ParentId: oldParent.ID}
	orphanChild := model.Menu{Name: "Menu Management", Path: "menu", Component: "/auth/menu"}
	assert.NoError(t, db.Create(&oldChild).Error)
	assert.NoError(t, db.Create(&orphanChild).Error)
	assert.NoError(t, db.Create(&model.RoleMenu{RoleId: role.ID, MenuId: orphanChild.ID}).Error)

	desiredParent := model.Menu{Name: "Quản lý phân quyền", Path: "/auth", Component: "Layout"}
	assert.NoError(t, upsertSeedParentMenu(db, &desiredParent))
	desiredChild := model.Menu{Name: "Quản lý menu", Path: "menu", Component: "/auth/menu", ParentId: desiredParent.ID}
	assert.NoError(t, upsertSeedChildMenu(db, &desiredChild))

	var parents, children []model.Menu
	assert.NoError(t, db.Where("parent_id = 0 AND path = ?", "/auth").Find(&parents).Error)
	assert.NoError(t, db.Where("path = ?", "menu").Find(&children).Error)
	assert.Len(t, parents, 1)
	assert.Len(t, children, 1)
	assert.Equal(t, "Quản lý phân quyền", parents[0].Name)
	assert.Equal(t, "Quản lý menu", children[0].Name)
	assert.Equal(t, parents[0].ID, children[0].ParentId)

	var bindingCount int64
	assert.NoError(t, db.Model(&model.RoleMenu{}).
		Where("role_id = ? AND menu_id = ?", role.ID, children[0].ID).
		Count(&bindingCount).Error)
	assert.Equal(t, int64(1), bindingCount)
}
