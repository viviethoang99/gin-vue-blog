package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Menu struct{}

type MenuTreeVO struct {
	model.Menu
	Children []MenuTreeVO `json:"children"`
}

// @Summary 获取当前用户菜单
// @Description 超级管理员返回全部菜单, 其他用户按角色返回
// @Tags Menu
// @Produce json
// @Success 0 {object} Response[[]MenuTreeVO]
// @Security ApiKeyAuth
// @Router /menu/user/list [get]
func (*Menu) GetUserMenu(c *gin.Context) {
	db := GetDB(c)
	auth, ok := MustCurrentUserAuth(c)
	if !ok {
		return
	}

	var menus []model.Menu
	var err error

	if auth.IsSuper {
		menus, err = model.GetAllMenuList(db)
	} else {
		menus, err = model.GetMenuListByUserId(GetDB(c), auth.ID)
	}

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menus))
}

// @Summary 获取菜单列表（树形）
// @Description 获取菜单列表, 支持关键字过滤
// @Tags Menu
// @Produce json
// @Param keyword query string false "关键字"
// @Success 0 {object} Response[[]MenuTreeVO]
// @Security ApiKeyAuth
// @Router /menu/list [get]
func (*Menu) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")

	menuList, _, err := model.GetMenuList(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menuList))
}

// @Summary 新增或编辑菜单
// @Description 新增或编辑菜单
// @Tags Menu
// @Accept json
// @Produce json
// @Param form body model.Menu true "新增或编辑菜单"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /menu [post]
func (*Menu) SaveOrUpdate(c *gin.Context) {
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := model.SaveOrUpdateMenu(GetDB(c), &req); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// @Summary 删除菜单
// @Description 被角色使用或存在子菜单时不允许删除
// @Tags Menu
// @Produce json
// @Param id path int true "菜单 ID"
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /menu/{id} [delete]
func (*Menu) Delete(c *gin.Context) {
	menuId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// Check whether the menu to delete is used by any role
	// 这里不能忽略 error: 查询失败时 use 是 false, 会把仍被角色引用的菜单删掉
	use, err := model.CheckMenuInUse(db, menuId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	if use {
		ReturnError(c, g.ErrMenuUsedByRole, nil)
		return
	}

	// If it's a first-level menu, check if it has child menus
	menu, err := model.GetMenuById(db, menuId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrMenuNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// Cannot delete if a first-level menu has child menus
	if menu.ParentId == 0 {
		has, err := model.CheckMenuHasChild(db, menuId)
		if err != nil {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
		if has {
			ReturnError(c, g.ErrMenuHasChildren, nil)
			return
		}
	}

	if err = model.DeleteMenu(db, menuId); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// @Summary 获取菜单选项（树形）
// @Description 用于角色分配菜单
// @Tags Menu
// @Produce json
// @Success 0 {object} Response[[]TreeOptionVO]
// @Security ApiKeyAuth
// @Router /menu/option [get]
func (*Menu) GetOption(c *gin.Context) {
	menus, _, err := model.GetMenuList(GetDB(c), "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	result := make([]TreeOptionVO, 0)
	for _, menu := range menus2MenuVos(menus) {
		option := TreeOptionVO{ID: menu.ID, Label: menu.Name}
		for _, child := range menu.Children {
			option.Children = append(option.Children, TreeOptionVO{ID: child.ID, Label: child.Name})
		}
		result = append(result, option)
	}

	ReturnSuccess(c, result)
}

// Build tree structure for menus, []Menu => []MenuVo
func menus2MenuVos(menus []model.Menu) []MenuTreeVO {
	result := make([]MenuTreeVO, 0)

	firstLevelMenus := getFirstLevelMenus(menus)
	childrenMap := getMenuChildrenMap(menus)

	for _, first := range firstLevelMenus {
		menu := MenuTreeVO{Menu: first}
		for _, childMenu := range childrenMap[first.ID] {
			menu.Children = append(menu.Children, MenuTreeVO{Menu: childMenu})
		}
		delete(childrenMap, first.ID)
		result = append(result, menu)
	}

	sortMenu(result)
	return result
}

// Filter first-level menus (parentId == 0)
func getFirstLevelMenus(menuList []model.Menu) []model.Menu {
	firstLevelMenus := make([]model.Menu, 0)
	for _, menu := range menuList {
		if menu.ParentId == 0 {
			firstLevelMenus = append(firstLevelMenus, menu)
		}
	}
	return firstLevelMenus
}

// key is menu ID, value is the corresponding child menu list
func getMenuChildrenMap(menus []model.Menu) map[int][]model.Menu {
	childrenMap := make(map[int][]model.Menu)
	for _, menu := range menus {
		if menu.ParentId != 0 {
			childrenMap[menu.ParentId] = append(childrenMap[menu.ParentId], menu)
		}
	}
	return childrenMap
}

// Sort by orderNum ascending (including child menus)
func sortMenu(menus []MenuTreeVO) {
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].OrderNum < menus[j].OrderNum
	})
	for i := range menus {
		sort.Slice(menus[i].Children, func(j, k int) bool {
			return menus[i].Children[j].OrderNum < menus[i].Children[k].OrderNum
		})
	}
}

// Build user menu tree structure, []Menu => []MenuVO
// func menus2UserMenuVos(menus []model.Menu) []MenuTreeVO {
// 	firstLevelMenuList := getFirstLevelMenus(menus)
// 	childrenMap := getMenuChildrenMap(menus)

// 	result := make([]MenuTreeVO, 0)

// 	// Iterate first-level menus and build UserMenu entries
// 	for _, firstLevelMenu := range firstLevelMenuList {
// 		var menuVO MenuTreeVO             // current user menu
// 		var userMenuChildren []MenuTreeVO // current user menu's children

// 		children := childrenMap[firstLevelMenu.ID] // child menus
// 		if len(children) > 0 {                     // has child menus
// 			menuVO = menu2UserMenuVo(firstLevelMenu) // [Menu] -> [UserMenu]
// 			// userMenu.Path = ""                           // TODO: Must outer path be ""?
// 			sortMenu(children) // sort child menus by OrderNum
// 			// Iterate child menus and build user menus
// 			for _, child := range children {
// 				userMenuChildren = append(userMenuChildren, menu2UserMenuVo(child))
// 			}
// 		} else { // No child menus: use first-level menu to construct a Layout user menu; put original menu as its child
// 			menuVO = MenuTreeVO{
// 				Menu: firstLevelMenu,
// 				// ID:        firstLevelMenu.ID,
// 				// Path:      firstLevelMenu.Path,
// 				// Name:      firstLevelMenu.Name,      // *
// 				// Component: firstLevelMenu.Component, // ! "Layout" ?
// 				// OrderNum:  firstLevelMenu.OrderNum,
// 				// IsHidden:  firstLevelMenu.Hidden,
// 				// KeepAlive: firstLevelMenu.KeepAlive,
// 				// Redirect:  firstLevelMenu.Redirect,
// 			}
// 			tmpUserMenu := menu2UserMenuVo(firstLevelMenu)
// 			// tmpUserMenu.Path = "" // TODO: consider this
// 			userMenuChildren = append(userMenuChildren, tmpUserMenu)
// 		}
// 		menuVO.Children = userMenuChildren
// 		result = append(result, menuVO)
// 	}
// 	return result
// }

/*
Menu data: []menuVo
{
	"id": 1,
	"name": "Home",
	"path": "/",
	"component": "/home/Home.vue",
	"icon": "el-icon-myshouye",
	"createTime": "2021-01-26T17:06:51",
	"orderNum": 1,
	"isDisable": null,
	"isHidden": 0,
	"children": []
},
{
	"id": 2,
	"name": "Article Management",
	"path": "/article-submenu",
	"component": "Layout",
	"icon": "el-icon-mywenzhang-copy",
	"createTime": "2021-01-25T20:43:07",
	"orderNum": 2,
	"isDisable": null,
	"isHidden": 0,
	"children": [
		{
			"id": 6,
			"name": "Publish Article",
			"path": "/articles",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"createTime": "2021-01-26T14:30:48",
			"orderNum": 1,
			"isDisable": null,
			"isHidden": 0,
			"children": null
		},
		{
			"id": 7,
			"name": "Edit Article",
			"path": "/articles/*",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"createTime": "2021-01-26T14:31:32",
			"orderNum": 2,
			"isDisable": null,
			"isHidden": 1,
			"children": null
		},
		{
			"id": 8,
			"name": "Article List",
			"path": "/article-list",
			"component": "/article/ArticleList.vue",
			"icon": "el-icon-mywenzhangliebiao",
			"createTime": "2021-01-26T14:32:13",
			"orderNum": 3,
			"isDisable": null,
			"isHidden": 0,
			"children": null
		}
	]
},
*/

/*
User menu data: []userMenuVo

{
	"name": null,
	"path": "/",
	"component": "Layout",
	"icon": null,
	"hidden": false,
	"children": [
		{
			"name": "Home",
			"path": "",
			"component": "/home/Home.vue",
			"icon": "el-icon-myshouye",
			"hidden": null,
			"children": null
		}
	]
},
{
	"name": "Article Management",
	"path": "/article-submenu",
	"component": "Layout",
	"icon": "el-icon-mywenzhang-copy",
	"hidden": false,
	"children": [
		{
			"name": "Publish Article",
			"path": "/articles",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"hidden": false,
			"children": null
		},
		{
			"name": "Edit Article",
			"path": "/articles/*",
			"component": "/article/Article.vue",
			"icon": "el-icon-myfabiaowenzhang",
			"hidden": true,
			"children": null
		},
		{
			"name": "Article List",
			"path": "/article-list",
			"component": "/article/ArticleList.vue",
			"icon": "el-icon-mywenzhangliebiao",
			"hidden": false,
			"children": null
		}
	]
}
*/
