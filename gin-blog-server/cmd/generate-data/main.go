package main

import (
	"flag"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("c", "../../config.yml", "Configuration file path")
	typeName := flag.String("t", "all", "Data type to initialize: config | auth | page | all")
	flag.Parse()

	// Read configuration file based on command line parameters, other variable initialization depends on the configuration file object
	conf := g.ReadConfig(*configPath)

	//! Handle sqlite3 database path
	conf.SQLite.Dsn = "../" + conf.SQLite.Dsn
	conf.Server.DbLogMode = "silent"

	db := ginblog.InitDatabase(conf)

	switch *typeName {
	case "config":
		generateDefaultConfigs(db)
	case "auth":
		generateDefaultAuths(db)
	case "page":
		generateDefaultPages(db)
	case "all":
		fallthrough
	default:
		generateDefaultConfigs(db)
		generateDefaultPages(db)
		generateDefaultAuths(db)
	}
}

// Generate authentication-related information: roles, users, resources, menus
func generateDefaultAuths(db *gorm.DB) {
	generateDefaultRolesAndUsers(db)
	generateDefaultResources(db)
	generateDefaultMenus(db)
}

// Generate default pages
func generateDefaultPages(db *gorm.DB) {
	slog.Info("-----Initialize blog pages start-----")

	pages := []model.Page{
		{Name: "Home", Label: "home", Cover: "https://cdn.hahacode.cn/page/home.jpg"},
		{Name: "Archive", Label: "archive", Cover: "https://cdn.hahacode.cn/page/archive.png"},
		{Name: "Category", Label: "category", Cover: "https://cdn.hahacode.cn/page/category.png"},
		{Name: "Tag", Label: "tag", Cover: "https://cdn.hahacode.cn/page/tag.png"},
		{Name: "Links", Label: "link", Cover: "https://cdn.hahacode.cn/page/link.jpg"},
		{Name: "About", Label: "about", Cover: "https://cdn.hahacode.cn/page/about.jpg"},
		{Name: "Message", Label: "message", Cover: "https://cdn.hahacode.cn/page/message.jpeg"},
		{Name: "User Center", Label: "user", Cover: "https://cdn.hahacode.cn/page/user.jpg"},
		{Name: "Album", Label: "album", Cover: "https://cdn.hahacode.cn/page/album.png"},
		{Name: "Error Page", Label: "404", Cover: "https://cdn.hahacode.cn/page/404.jpg"},
		{Name: "Article List", Label: "article_list", Cover: "https://cdn.hahacode.cn/page/article_list.jpg"},
	}

	for _, page := range pages {
		if err := db.Create(&page).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(page.Name + " page data already exists")
			} else {
				slog.Error(page.Name + " page initialization failed" + err.Error())
			}
		}
	}

	slog.Info("-----Initialize blog pages end-----")
}

// Generate default configuration information
func generateDefaultConfigs(db *gorm.DB) {
	slog.Info("-----Initialize blog configuration start-----")

	configs := []model.Config{
		{Key: "website_avatar", Value: "https://foruda.gitee.com/avatar/1677041571085433939/5221991_szluyu99_1614389421.png", Desc: "Website Avatar"},
		{Key: "website_name", Value: "Zhenyu's Personal Blog", Desc: "Website Name"},
		{Key: "website_author", Value: "Zhenyu", Desc: "Website Author"},
		{Key: "website_intro", Value: "Let the past go with the wind", Desc: "Website Introduction"},
		{Key: "website_notice", Value: "Welcome to Zhenyu's personal blog, the project is still under development...", Desc: "Website Notice"},
		{Key: "website_createtime", Value: time.Now().Format(time.DateTime), Desc: "Website Creation Date"},
		{Key: "website_record", Value: "ICP Registration No. 2021032312", Desc: "Website Registration Number"},
		{Key: "qq", Value: "123456789", Desc: "QQ"},
		{Key: "github", Value: "https://github.com/szluyu99", Desc: "github"},
		{Key: "gitee", Value: "https://gitee.com/szluyu99", Desc: "gitee"},
		{Key: "tourist_avatar", Value: "https://cdn.hahacode.cn/config/tourist_avatar.png", Desc: "Default Tourist Avatar"},
		{Key: "user_avatar", Value: "https://cdn.hahacode.cn/config/user_avatar.png", Desc: "Default User Avatar"},
		{Key: "article_cover", Value: "https://cdn.hahacode.cn/config/default_article_cover.png", Desc: "Default Article Cover"},
		{Key: "is_comment_review", Value: "true", Desc: "Comment Default Review"},
		{Key: "is_message_review", Value: "true", Desc: "Message Default Review"},
	}

	for _, config := range configs {
		if err := db.Create(&config).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(config.Key + " configuration already exists")
			} else {
				slog.Error(config.Key + " configuration initialization failed" + err.Error())
			}
		}
	}

	slog.Info("-----Initialize blog configuration end-----")
}

// Generate 2 default roles and authentication information: admin, guest
func generateDefaultRolesAndUsers(db *gorm.DB) {
	slog.Info("-----Initialize default roles and users start-----")

	roles := []model.Role{
		{Name: "admin", Label: "Administrator"},
		{Name: "guest", Label: "Guest"},
	}

	for i := range roles {
		if err := db.Create(&roles[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(roles[i].Name + " role already exists")
			} else {
				slog.Error(roles[i].Name + " role initialization failed" + err.Error())
			}
		}
	}

	pwd, _ := utils.BcryptHash("123456")
	auths := []model.UserAuth{
		{
			Username: "admin",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "admin",
				Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
			},
		},
		{
			Username: "guest",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "guest",
				Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
			},
		},
	}

	for i := range auths {
		if err := db.Create(&auths[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(auths[i].Username + " user already exists")
			} else {
				slog.Error(auths[i].Username + " user initialization failed" + err.Error())
			}
		}
		// Create user role association
		db.Create(&model.UserAuthRole{UserAuthId: auths[i].ID, RoleId: roles[i].ID})
	}

	slog.Info("-----Initialize default roles and users end-----")
}

// Generate default interface resources
func generateDefaultResources(db *gorm.DB) {
	slog.Info("-----Initialize interface resources start-----")

	parents := []model.Resource{
		{Name: "Article Module"},
		{Name: "Category Module"},
		{Name: "Tag Module"},
		{Name: "Page Module"},
		{Name: "Link Module"},
		{Name: "Menu Module"},
		{Name: "Role Module"},
		{Name: "Resource Module"},
		{Name: "Comment Module"},
		{Name: "Message Module"},
		{Name: "File Module"},
		{Name: "Blog Info Module"},
		{Name: "User Info Module"},
		{Name: "Operation Log Module"},
	}
	for i := range parents {
		if err := db.Create(&parents[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(parents[i].Name + " resource already exists")
			} else {
				slog.Error(parents[i].Name + " resource initialization failed" + err.Error())
			}
		}
	}

	resources := []model.Resource{
		// Article Module
		{Name: "Article List", ParentId: parents[0].ID, Url: "/article/list", Method: "GET"},
		{Name: "Article Details", ParentId: parents[0].ID, Url: "/article/:id", Method: "GET"},
		{Name: "Add/Edit Article", ParentId: parents[0].ID, Url: "/article", Method: "POST"},
		{Name: "Update Article Soft Delete", ParentId: parents[0].ID, Url: "/article/soft-delete", Method: "PUT"},
		{Name: "Delete Article", ParentId: parents[0].ID, Url: "/article", Method: "DELETE"},
		{Name: "Modify Article Top", ParentId: parents[0].ID, Url: "/article/top", Method: "PUT"},
		{Name: "Export Article", ParentId: parents[0].ID, Url: "/article/export", Method: "POST"},
		{Name: "Import Article", ParentId: parents[0].ID, Url: "/article/import", Method: "POST"},
		// Category Module
		{Name: "Category List", ParentId: parents[1].ID, Url: "/category/list", Method: "GET"},
		{Name: "Add/Edit Category", ParentId: parents[1].ID, Url: "/category", Method: "POST"},
		{Name: "Delete Category", ParentId: parents[1].ID, Url: "/category", Method: "DELETE"},
		{Name: "Category Option List", ParentId: parents[1].ID, Url: "/category/option", Method: "GET"},
		// Tag Module
		{Name: "Tag List", ParentId: parents[2].ID, Url: "/tag/list", Method: "GET"},
		{Name: "Add/Edit Tag", ParentId: parents[2].ID, Url: "/tag", Method: "POST"},
		{Name: "Delete Tag", ParentId: parents[2].ID, Url: "/tag", Method: "DELETE"},
		{Name: "Tag Option List", ParentId: parents[2].ID, Url: "/tag/option", Method: "GET"},
		// Page Module
		{Name: "Page List", ParentId: parents[3].ID, Url: "/page/list", Method: "GET"},
		{Name: "Add/Edit Page", ParentId: parents[3].ID, Url: "/page", Method: "POST"},
		{Name: "Delete Page", ParentId: parents[3].ID, Url: "/page", Method: "DELETE"},
		// Link Module
		{Name: "Link List", ParentId: parents[4].ID, Url: "/link/list", Method: "GET"},
		{Name: "Add/Edit Link", ParentId: parents[4].ID, Url: "/link", Method: "POST"},
		{Name: "Delete Link", ParentId: parents[4].ID, Url: "/link", Method: "DELETE"},
		// Menu Module
		{Name: "Menu List", ParentId: parents[5].ID, Url: "/menu/list", Method: "GET"},
		{Name: "Add/Edit Menu", ParentId: parents[5].ID, Url: "/menu", Method: "POST"},
		{Name: "Delete Menu", ParentId: parents[5].ID, Url: "/menu", Method: "DELETE"},
		{Name: "Menu Option List (Tree)", ParentId: parents[5].ID, Url: "/menu/option", Method: "GET"},
		{Name: "Get Current User Menu", ParentId: parents[5].ID, Url: "/menu/user/list", Method: "GET"},
		// Role Module
		{Name: "Role List", ParentId: parents[6].ID, Url: "/role/list", Method: "GET"},
		{Name: "Add/Edit Role", ParentId: parents[6].ID, Url: "/role", Method: "POST"},
		{Name: "Delete Role", ParentId: parents[6].ID, Url: "/role", Method: "DELETE"},
		{Name: "Role Option List", ParentId: parents[6].ID, Url: "/role/option", Method: "GET"},
		// Resource Module
		{Name: "Resource List", ParentId: parents[7].ID, Url: "/resource/list", Method: "GET"},
		{Name: "Add/Edit Resource", ParentId: parents[7].ID, Url: "/resource", Method: "POST"},
		{Name: "Delete Resource", ParentId: parents[7].ID, Url: "/resource", Method: "DELETE"},
		{Name: "Resource Option List (Tree)", ParentId: parents[7].ID, Url: "/resource/option", Method: "GET"},
		{Name: "Modify Resource Anonymous Access", ParentId: parents[7].ID, Url: "/resource/anonymous", Method: "PUT"},
		// Comment Module
		{Name: "Comment List", ParentId: parents[8].ID, Url: "/comment/list", Method: "GET"},
		{Name: "Delete Comment", ParentId: parents[8].ID, Url: "/comment", Method: "DELETE"},
		{Name: "Modify Comment Review", ParentId: parents[8].ID, Url: "/comment/review", Method: "PUT"},
		// Message Module
		{Name: "Message List", ParentId: parents[9].ID, Url: "/message/list", Method: "GET"},
		{Name: "Delete Message", ParentId: parents[9].ID, Url: "/message", Method: "DELETE"},
		{Name: "Modify Message Review", ParentId: parents[9].ID, Url: "/message/review", Method: "PUT"},
		// File Module
		{Name: "File Upload", ParentId: parents[10].ID, Url: "/upload", Method: "POST"},
		// Blog Info Module
		{Name: "Get Blog Settings", ParentId: parents[11].ID, Url: "/setting/blog-config", Method: "GET"},
		{Name: "Get About Me", ParentId: parents[11].ID, Url: "/setting/about", Method: "GET"},
		{Name: "Modify Blog Settings", ParentId: parents[11].ID, Url: "/setting/blog-config", Method: "PUT"},
		{Name: "Modify About Me", ParentId: parents[11].ID, Url: "/setting/about", Method: "PUT"},
		{Name: "Get Backend Home Info", ParentId: parents[11].ID, Url: "/home", Method: "GET"},
		// User Info Module
		{Name: "User List", ParentId: parents[12].ID, Url: "/user/list", Method: "GET"},
		{Name: "Get Current User Info", ParentId: parents[12].ID, Url: "/user/info", Method: "GET"},
		{Name: "Modify User Info", ParentId: parents[12].ID, Url: "/user", Method: "PUT"},
		{Name: "Get Online User List", ParentId: parents[12].ID, Url: "/user/online", Method: "GET"},
		{Name: "Force User Offline", ParentId: parents[12].ID, Url: "/user/offline", Method: "DELETE"},
		{Name: "Modify Current User Password", ParentId: parents[12].ID, Url: "/user/current/password", Method: "PUT"},
		{Name: "Modify Current User Info", ParentId: parents[12].ID, Url: "/user/current", Method: "PUT"},
		{Name: "Modify User Disable", ParentId: parents[12].ID, Url: "/user/disable", Method: "PUT"},
		// Operation Log Module
		{Name: "Log List", ParentId: parents[13].ID, Url: "/operation/log/list", Method: "GET"},
		{Name: "Delete Operation Log", ParentId: parents[13].ID, Url: "/operation/log", Method: "DELETE"},
	}

	for i := range resources {
		if err := db.Create(&resources[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(resources[i].Name + " resource already exists")
			} else {
				slog.Error(resources[i].Name + " resource initialization failed" + err.Error())
			}
		}
	}

	// Load all resources
	db.Find(&resources)

	// Add all resource access permissions to admin role
	var adminRole model.Role
	if err := db.Where("name", "admin").First(&adminRole).Error; err == nil {
		for _, resource := range resources {
			if resource.ID != 0 {
				if err := db.Create(&model.RoleResource{RoleId: adminRole.ID, ResourceId: resource.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("admin role menu association initialization failed" + err.Error())
					} else {
						slog.Error("admin role menu association initialization failed" + err.Error())
					}
				}
			}
		}
	}

	// Add query resource access permissions to guest
	var guestRole model.Role
	if err := db.Where("name", "guest").First(&guestRole).Error; err == nil {
		for _, resource := range resources {
			if resource.ID != 0 && resource.Method == "GET" {
				if err := db.Create(&model.RoleResource{RoleId: guestRole.ID, ResourceId: resource.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("guest role menu association initialization failed" + err.Error())
					} else {
						slog.Error("guest role menu association initialization failed" + err.Error())
					}
				}
			}
		}
	}

	slog.Info("-----Initialize interface resources end-----")
}

// Generate default menus
func generateDefaultMenus(db *gorm.DB) {
	slog.Info("-----Initialize menus start-----")

	parents := []model.Menu{
		{Name: "Home", Path: "/home", Icon: "ic:sharp-home", OrderNum: 0, Component: "/home", Redirect: "/home", Catalogue: true}, // catalogue
		{Name: "Article Management", Path: "/article", Icon: "ic:twotone-article", OrderNum: 1, Component: "Layout", Redirect: "/article/list"},
		{Name: "Permission Management", Path: "/auth", Icon: "cib:adguard", OrderNum: 3, Component: "Layout", Redirect: "/auth/menu"},
		{Name: "Message Management", Path: "/message", Icon: "ic:twotone-email", OrderNum: 2, Component: "Layout", Redirect: "/message/comment"},
		{Name: "User Management", Path: "/user", Icon: "ph:user-list-bold", OrderNum: 4, Component: "Layout", Redirect: "/user/list"},
		{Name: "Log Management", Path: "/log", Icon: "material-symbols:receipt-long-outline-rounded", OrderNum: 6, Component: "Layout", Redirect: "/log/operation"},
		{Name: "System Management", Path: "/setting", Icon: "ion:md-settings", OrderNum: 5, Component: "Layout", Redirect: "/setting/website"},
		{Name: "Personal Center", Path: "/profile", Icon: "mdi:account", OrderNum: 7, Component: "/profile", Redirect: "/profile", Catalogue: true}, // catalogue
	}

	for i := range parents {
		if err := db.Create(&parents[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(parents[i].Name + " menu already exists")
			} else {
				slog.Error(parents[i].Name + " menu initialization failed" + err.Error())
			}
		}
	}

	menus := []model.Menu{
		// Article Management
		{Name: "Publish Article", Path: "write", Component: "/article/write", Icon: "icon-park-outline:write", OrderNum: 1, ParentId: parents[1].ID},
		{Name: "Article List", Path: "list", Component: "/article/list", Icon: "material-symbols:format-list-bulleted", OrderNum: 2, ParentId: parents[1].ID},
		{Name: "Category Management", Path: "category", Component: "/article/category", Icon: "tabler:category", OrderNum: 3, ParentId: parents[1].ID},
		{Name: "Tag Management", Path: "tag", Component: "/article/tag", Icon: "tabler:tag", OrderNum: 4, ParentId: parents[1].ID},
		{Name: "Edit Article", Path: "write/:id", Component: "/article/write", Icon: "icon-park-outline:write", OrderNum: 1, ParentId: parents[1].ID, Hidden: true},
		// Permission Management
		{Name: "Menu Management", Path: "menu", Component: "/auth/menu", Icon: "ic:twotone-menu-book", OrderNum: 1, ParentId: parents[2].ID},
		{Name: "Interface Management", Path: "resource", Component: "/auth/resource", Icon: "mdi:api", OrderNum: 2, ParentId: parents[2].ID},
		{Name: "Role Management", Path: "role", Component: "/auth/role", Icon: "carbon:user-role", OrderNum: 3, ParentId: parents[2].ID},
		// Message Management
		{Name: "Comment Management", Path: "comment", Component: "/message/comment", Icon: "ic:twotone-comment", OrderNum: 1, ParentId: parents[3].ID},
		{Name: "Message Management", Path: "leave-msg", Component: "/message/leave-msg", Icon: "ic:twotone-message", OrderNum: 2, ParentId: parents[3].ID},
		// User Management
		{Name: "User List", Path: "list", Component: "/user/list", Icon: "mdi:account", OrderNum: 1, ParentId: parents[4].ID},
		{Name: "Online Users", Path: "online", Component: "/user/online", Icon: "ic:outline-online-prediction", OrderNum: 2, ParentId: parents[4].ID},
		// Log Management
		{Name: "Operation Log", Path: "operation", Component: "/log/operation", Icon: "mdi:book-open-page-variant-outline", OrderNum: 1, ParentId: parents[5].ID},
		{Name: "Login Log", Path: "login", Component: "/log/login", Icon: "material-symbols:login", OrderNum: 2, ParentId: parents[5].ID},
		// System Management
		{Name: "Website Management", Path: "website", Component: "/setting/website", Icon: "el:website", OrderNum: 1, ParentId: parents[6].ID},
		{Name: "Page Management", Path: "page", Component: "/setting/page", Icon: "iconoir:journal-page", OrderNum: 2, ParentId: parents[6].ID},
		{Name: "Link Management", Path: "link", Component: "/setting/link", Icon: "mdi:telegram", OrderNum: 3, ParentId: parents[6].ID},
		{Name: "About Me", Path: "about", Component: "/setting/about", Icon: "cib:about-me", OrderNum: 4, ParentId: parents[6].ID},
	}

	for i := range menus {
		if err := db.Create(&menus[i]).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
				slog.Info(menus[i].Name + " menu already exists")
			} else {
				slog.Error(menus[i].Name + " menu initialization failed" + err.Error())
			}
		}
	}

	// Load all menus
	db.Find(&menus)

	// Add all menu access permissions to admin role
	var adminRole model.Role
	if err := db.Where("name", "admin").First(&adminRole).Error; err == nil {
		for _, menu := range menus {
			if menu.ID != 0 {
				if err := db.Create(&model.RoleMenu{RoleId: adminRole.ID, MenuId: menu.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("admin role menu association initialization failed" + err.Error())
					} else {
						slog.Error("admin role menu association initialization failed" + err.Error())
					}
				}
			}
		}
	}

	// Add all menu access permissions to guest role
	var guestRole model.Role
	if err := db.Where("name", "guest").First(&guestRole).Error; err == nil {
		for _, menu := range menus {
			if menu.ID != 0 {
				if err := db.Create(&model.RoleMenu{RoleId: guestRole.ID, MenuId: menu.ID}).Error; err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
						slog.Info("guest role menu association initialization failed" + err.Error())
					} else {
						slog.Error("guest role menu association initialization failed" + err.Error())
					}
				}
			}
		}
	}

	slog.Info("-----Initialize menus end-----")
}
