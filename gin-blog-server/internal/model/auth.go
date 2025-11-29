package model

import (
	"encoding/json"
	"gin-blog/internal/utils"
	"log/slog"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Access control: 7 tables (4 models + 3 relations)

type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:Login type" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:Login IP address" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP source" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // Super admin can only be set in the backend

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"`
	Roles      []*Role   `json:"roles" gorm:"many2many:user_auth_role"`
}

func (u *UserAuth) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources []Resource `json:"resources" gorm:"many2many:role_resource"`
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menu"`
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_role"`
}

type Resource struct {
	Model
	Name      string `gorm:"unique;type:varchar(200)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255)" json:"url"`
	Method    string `gorm:"type:varchar(10)" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles []*Role `json:"roles" gorm:"many2many:role_resource"`
}

/*
Menu design:

Catalogue: catalogue === true
	- If it is a catalogue, it appears as a single item and does not expand a submenu (e.g., "Home", "Profile").
	- If not a catalogue and parent_id is 0, it is a first-level menu that can expand submenus (e.g., under "Article Management" there are "Article List", "Article Category", "Article Tag").
	- If not a catalogue and parent_id is not 0, it is a second-level menu.

Hidden: hidden
	- If hidden, it does not appear in the sidebar menu.

External: external, external_link
	- If external, clicking opens in a new window.
*/
type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`
		Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(200)" json:"name"` // Menu name
		Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"`  // Route path
		Component    string `gorm:"type:varchar(50)" json:"component"`                           // Component path
		Icon         string `gorm:"type:varchar(50)" json:"icon"`                                // Icon
		OrderNum     int8   `json:"order_num"`                                                   // Order
		Redirect     string `gorm:"type:varchar(50)" json:"redirect"`                            // Redirect URL
		Catalogue    bool   `json:"is_catalogue"`                                                // Is catalogue
		Hidden       bool   `json:"is_hidden"`                                                   // Hidden
		KeepAlive    bool   `json:"keep_alive"`                                                  // Keep alive (cache)
		External     bool   `json:"is_external"`                                                 // External link
		ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`                      // External URL

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"`
}

type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}

type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	MenuId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}

type RoleVO struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	IsDisable   bool      `json:"is_disable"`
	ResourceIds []int     `json:"resource_ids" gorm:"-"`
	MenuIds     []int     `json:"menu_ids" gorm:"-"`
}

// Menu

func SaveOrUpdateMenu(db *gorm.DB, menu *Menu) error {
	var result *gorm.DB

	if menu.ID > 0 {
		result = db.Model(menu).
			Select("name", "path", "component", "icon", "redirect", "parent_id", "order_num", "catalogue", "hidden", "keep_alive", "external").
			Updates(menu)
	} else {
		result = db.Create(menu)
	}

	return result.Error
}

func GetMenuIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&RoleMenu{}).Where("role_id = ?", roleId).Pluck("menu_id", &ids)
	return ids, result.Error
}

func GetMenuById(db *gorm.DB, id int) (menu *Menu, err error) {
	result := db.First(&menu, id)
	return menu, result.Error
}

func CheckMenuInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&RoleMenu{}).Where("menu_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func CheckMenuHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&Menu{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Get all menu list (for super admin)
func GetAllMenuList(db *gorm.DB) (menu []Menu, err error) {
	result := db.Find(&menu)
	return menu, result.Error
}

// Get menu list by user_id
func GetMenuListByUserId(db *gorm.DB, id int) (menus []Menu, err error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Model: Model{ID: id}}).
		Preload("Roles").Preload("Roles.Menus").
		First(&userAuth)

	if result.Error != nil {
		return nil, result.Error
	}

	set := make(map[int]Menu)
	for _, role := range userAuth.Roles {
		for _, menu := range role.Menus {
			set[menu.ID] = menu
		}
	}

	for _, menu := range set {
		menus = append(menus, menu)
	}

	return menus, nil
}

func GetMenuList(db *gorm.DB, keyword string) (list []Menu, total int64, err error) {
	db = db.Model(&Menu{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	result := db.Count(&total).Find(&list)
	return list, total, result.Error
}

func DeleteMenu(db *gorm.DB, id int) error {
	result := db.Delete(&Menu{}, id)
	return result.Error
}

// Resource

func SaveOrUpdateResource(db *gorm.DB, id, pid int, name, url, method string) error {
	resource := Resource{
		Model:    Model{ID: id},
		Name:     name,
		Url:      url,
		Method:   method,
		ParentId: pid,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&resource)
	} else {
		result = db.Create(&resource)
		// TODO: Front-end workaround
		// - Fix a front-end bug: after cascade-selecting a parent node, a newly added child node appears selected by default though it isn't actually selected.
		// - Workaround: After adding a child node, remove the association between its parent node and roles.
		// dao.Delete(model.RoleResource{}, "resource_id", data.ParentId)
	}
	return result.Error
}

func GetResourceIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&RoleResource{}).
		Where("role_id = ?", roleId).
		Pluck("resource_id", &ids)
	return ids, result.Error
}

func GetResourceList(db *gorm.DB, keyword string) (list []Resource, err error) {
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	result := db.Find(&list)
	return list, result.Error
}

func GetResourceListByIds(db *gorm.DB, ids []int) (list []Resource, err error) {
	result := db.Where("id in ?", ids).Find(&list)
	return list, result.Error
}

// Role

func SaveOrUpdateRole(db *gorm.DB, id int, name, label string, isDisable bool) error {
	role := Role{
		Model:     Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&role)
	} else {
		result = db.Create(&role)
	}

	return result.Error
}

func GetRoleOption(db *gorm.DB) (list []OptionVO, err error) {
	result := db.Model(&Role{}).Select("id", "name").Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func GetRoleList(db *gorm.DB, num, size int, keyword string) (list []RoleVO, total int64, err error) {
	db = db.Model(&Role{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Select("id", "name", "label", "created_at", "is_disable").
		Scopes(Paginate(num, size)).
		Find(&list)
	return list, total, result.Error
}

func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int, err error) {
	result := db.
		Model(&UserAuthRole{UserAuthId: userAuthId}).
		Pluck("role_id", &ids)
	return ids, result.Error
}

func SaveRole(db *gorm.DB, name, label string) error {
	role := Role{
		Name:  name,
		Label: label,
	}
	result := db.Create(&role)
	return result.Error
}

func UpdateRole(db *gorm.DB, id int, name, label string, isDisable bool, resourceIds, menuIds []int) error {
	role := Role{
		Model:     Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := db.Model(&role).Select("name", "label", "is_disable").Updates(&role).Error; err != nil {
			return err
		}

		// role_resource
		if err := db.Delete(&RoleResource{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		for _, rid := range resourceIds {
			if err := db.Create(&RoleResource{RoleId: role.ID, ResourceId: rid}).Error; err != nil {
				return err
			}
		}

		// role_menu
		if err := db.Delete(&RoleMenu{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		for _, mid := range menuIds {
			if err := db.Create(&RoleMenu{RoleId: role.ID, MenuId: mid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete roles: transactionally delete role, role_resource, and role_menu
func DeleteRoles(db *gorm.DB, ids []int) error {
	return db.Transaction(func(tx *gorm.DB) error {

		result := db.Delete(&Role{}, "id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&RoleResource{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&RoleMenu{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

// UserAuth

func GetUserAuthInfoById(db *gorm.DB, id int) (*UserAuth, error) {
	var userAuth = UserAuth{Model: Model{ID: id}}
	result := db.Model(&userAuth).
		Preload("Roles").Preload("UserInfo").
		First(&userAuth)
	return &userAuth, result.Error
}

// Register a new user
func CreateNewUser(db *gorm.DB, username, password string) (*UserAuth, *UserInfo, *UserAuthRole, error) {
	// Create user info
	num, err := Count(db, &UserInfo{})
	if err != nil {
		slog.Info(err.Error())
	}
	number := strconv.Itoa(num)
	userinfo := &UserInfo{
		Email:    username,
		Nickname: "Visitor" + number,
		Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
		Intro:    "I am user #" + number + " of this application",
	}
	result := db.Create(&userinfo)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	// Create user auth first
	pass, _ := utils.BcryptHash(password)
	userauth := &UserAuth{
		Username:   username,
		Password:   pass,
		UserInfoId: userinfo.ID,
	}

	result = db.Create(&userauth)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	// Then create role association record
	user_role := &UserAuthRole{
		UserAuthId: userauth.ID,
		RoleId:     2, // Default role is visitor
	}
	result = db.Create(&user_role)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	return userauth, userinfo, user_role, result.Error
}
