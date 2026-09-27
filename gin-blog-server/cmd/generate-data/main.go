package main

import (
	"errors"
	"flag"
	"fmt"
	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 种子数据里的图片: 原来指向 cdn.hahacode.cn, 该图床已经失效(连接超时),
// 换成本仓库 images/ 目录下的图片, 由 GitHub 直接提供
const imgBase = "https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images"

func main() {
	configPath := flag.String("c", "../../config.yml", "Configuration file path")
	typeName := flag.String("t", "all", "要初始化的数据类型: config | auth | page | demo | all")
	// 从 cmd/generate-data 目录运行时, sqlite 文件在上一级(与 server 的工作目录 cmd/ 一致);
	// 容器里二进制和配置文件同级, 需要显式关掉
	sqliteParent := flag.Bool("sqlite-parent", true, "sqlite 数据库文件是否在上级目录")
	flag.Parse()

	// Read configuration file based on command line parameters, other variable initialization depends on the configuration file object
	conf := g.ReadConfig(*configPath)

	if *sqliteParent {
		conf.SQLite.Dsn = "../" + conf.SQLite.Dsn
	}
	conf.Server.DbLogMode = "silent"

	db := ginblog.InitDatabase(conf)

	switch *typeName {
	case "config":
		generateDefaultConfigs(db)
	case "auth":
		generateDefaultAuths(db)
	case "page":
		generateDefaultPages(db)
	case "demo":
		generateDemoContent(db)
	case "all":
		fallthrough
	default:
		generateDefaultConfigs(db)
		generateDefaultPages(db)
		generateDefaultAuths(db)
	}
}

// 生成样例内容: 分类, 标签, 文章, 评论, 留言, 友链, 说说
//
// 只在库里一篇文章都没有时才灌, 所以可以重复执行。
// 内容数据只是本地测试用, 不放进 all 里, 需要时显式 -t demo。
func generateDemoContent(db *gorm.DB) {
	slog.Info("-----初始化样例内容 start-----")
	if err := model.SeedDemoContent(db); err != nil {
		slog.Error("样例内容初始化失败: " + err.Error())
	}
	slog.Info("-----初始化样例内容 end-----")
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
		{Name: "Home", Label: "home", Cover: imgBase + "/page/home.jpg"},
		{Name: "Archive", Label: "archive", Cover: imgBase + "/page/archive.png"},
		{Name: "Category", Label: "category", Cover: imgBase + "/page/category.png"},
		{Name: "Tag", Label: "tag", Cover: imgBase + "/page/tag.png"},
		{Name: "Links", Label: "link", Cover: imgBase + "/page/link.jpg"},
		{Name: "About", Label: "about", Cover: imgBase + "/page/about.jpg"},
		{Name: "Message", Label: "message", Cover: imgBase + "/page/message.jpeg"},
		{Name: "Personal Center", Label: "user", Cover: imgBase + "/page/user.jpg"},
		{Name: "Album", Label: "album", Cover: imgBase + "/page/album.png"},
		{Name: "说说", Label: "talk", Cover: imgBase + "/page/talking.jpg"},
		{Name: "Error Page", Label: "404", Cover: imgBase + "/page/404.jpg"},
		{Name: "Article List", Label: "article_list", Cover: imgBase + "/page/article_list.jpg"},
	}

	for _, page := range pages {
		if err := db.Create(&page).Error; err != nil {
			if isDuplicate(err) {
				slog.Debug(page.Name + " page data already exists")
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
		{Key: "website_avatar", Value: imgBase + "/common/header.jpeg", Desc: "Website Avatar"},
		{Key: "website_name", Value: "Zhenyu's Personal Blog", Desc: "Website Name"},
		{Key: "website_author", Value: "Zhenyu", Desc: "Website Author"},
		{Key: "website_intro", Value: "Let the past go with the wind", Desc: "Website Introduction"},
		{Key: "website_notice", Value: "Welcome to Zhenyu's personal blog, the project is still under development...", Desc: "Website Notice"},
		{Key: "website_createtime", Value: time.Now().Format(time.DateTime), Desc: "Website Creation Date"},
		{Key: "website_record", Value: "ICP Registration No. 2021032312", Desc: "Website Registration Number"},
		{Key: "qq", Value: "123456789", Desc: "QQ"},
		{Key: "github", Value: "https://github.com/szluyu99", Desc: "github"},
		{Key: "gitee", Value: "https://gitee.com/szluyu99", Desc: "gitee"},
		{Key: "tourist_avatar", Value: imgBase + "/config/tourist_avatar.jpeg", Desc: "Default Tourist Avatar"},
		{Key: "user_avatar", Value: imgBase + "/config/user_avatar.jpeg", Desc: "Default User Avatar"},
		{Key: "article_cover", Value: imgBase + "/config/default_article_cover.png", Desc: "Default Article Cover"},
		// 名字读起来像「需要审核」, 实际语义是「免审核」: true = 新内容直接展示,
		// false = 要在后台点「通过」。Desc 写清楚, 免得照名字理解反了
		{Key: "is_comment_review", Value: "true", Desc: "评论免审核(true 新评论直接展示, false 需后台通过)"},
		{Key: "is_message_review", Value: "true", Desc: "留言免审核(true 新留言直接展示, false 需后台通过)"},
	}

	for _, config := range configs {
		if err := db.Create(&config).Error; err != nil {
			if isDuplicate(err) {
				slog.Debug(config.Key + " configuration already exists")
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
			if isDuplicate(err) {
				slog.Debug(roles[i].Name + " role already exists")
				// 取回已有 ID, 否则后面的关联关系会写成 0
				db.Where("name", roles[i].Name).First(&roles[i])
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
				Avatar:   imgBase + "/config/user_avatar.jpeg",
			},
		},
		{
			Username: "guest",
			Password: pwd,
			UserInfo: &model.UserInfo{
				Nickname: "guest",
				Avatar:   imgBase + "/config/user_avatar.jpeg",
			},
		},
	}

	for i := range auths {
		if err := db.Create(&auths[i]).Error; err != nil {
			if isDuplicate(err) {
				slog.Debug(auths[i].Username + " user already exists")
				// 取回已有 ID, 否则下面会插入 user_auth_id = 0 的脏数据
				db.Where("username", auths[i].Username).First(&auths[i])
			} else {
				slog.Error(auths[i].Username + " user initialization failed" + err.Error())
			}
		}
		// Create user role association
		if auths[i].ID != 0 && roles[i].ID != 0 {
			db.Create(&model.UserAuthRole{UserAuthId: auths[i].ID, RoleId: roles[i].ID})
		}
	}

	slog.Info("-----Initialize default roles and users end-----")
}

// Generate default interface resources
//
// 资源定义见 internal/model/seed_resource.go, 与后台路由一一对应。
// 这里是对账而不是只增不删: 接口下线或改名后, 旧资源如果留在表里会继续
// 挂在角色上, 之后路由被复用时权限就凭空对上了。
func generateDefaultResources(db *gorm.DB) {
	slog.Info("-----Initialize interface resources start-----")

	// 期望存在的接口资源, key 为 "METHOD URL"
	wanted := make(map[string]model.Resource)
	// 期望存在的模块(父资源)名称
	moduleNames := make(map[string]bool)

	for _, module := range model.AdminResources {
		moduleNames[module.Name] = true

		parent := model.Resource{Name: module.Name}
		err := db.Where("name", module.Name).First(&parent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&parent).Error; err != nil {
				slog.Error(module.Name + " resource initialization failed" + err.Error())
				continue
			}
		} else if err != nil {
			slog.Error(module.Name + " 资源查询失败" + err.Error())
			continue
		}

		for _, item := range module.Items {
			wanted[item.Method+" "+item.Url] = model.Resource{
				Name:     item.Name,
				ParentId: parent.ID,
				Url:      item.Url,
				Method:   item.Method,
			}
		}
	}

	// 新增或更新: 以 url + method 定位, 名称和所属模块允许变更
	for _, want := range wanted {
		var exist model.Resource
		err := db.Where(&model.Resource{Url: want.Url, Method: want.Method}).First(&exist).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&want).Error; err != nil {
				slog.Error(want.Name + " resource initialization failed" + err.Error())
			}
			continue
		}
		if err != nil {
			slog.Error(want.Name + " 资源查询失败" + err.Error())
			continue
		}
		if exist.Name != want.Name || exist.ParentId != want.ParentId {
			if err := db.Model(&exist).
				Updates(map[string]any{"name": want.Name, "parent_id": want.ParentId}).Error; err != nil {
				slog.Error(want.Name + " 资源更新失败" + err.Error())
			}
		}
	}

	// 清理代码中已经不存在的资源, 连同它的角色关联
	// 注意只删资源本身, 不动管理员在后台手动给角色配的其他权限
	var all []model.Resource
	if err := db.Find(&all).Error; err != nil {
		slog.Error("资源列表查询失败" + err.Error())
		return
	}

	var stale []int
	for _, r := range all {
		if r.Url == "" && r.Method == "" { // 模块(父资源)
			if !moduleNames[r.Name] {
				stale = append(stale, r.ID)
			}
			continue
		}
		if _, ok := wanted[r.Method+" "+r.Url]; !ok {
			stale = append(stale, r.ID)
		}
	}

	if len(stale) > 0 {
		if err := db.Delete(&model.RoleResource{}, "resource_id in ?", stale).Error; err != nil {
			slog.Error("清理过期资源的角色关联失败" + err.Error())
		}
		if err := db.Delete(&model.Resource{}, "id in ?", stale).Error; err != nil {
			slog.Error("清理过期资源失败" + err.Error())
		}
		slog.Info(fmt.Sprintf("清理了 %d 条代码中已不存在的资源", len(stale)))
	}

	// 重新加载, 下面按最新的资源表建立角色关联
	if err := db.Find(&all).Error; err != nil {
		slog.Error("资源列表查询失败" + err.Error())
		return
	}

	// Add all resource access permissions to admin role
	var adminRole model.Role
	if err := db.Where("name", "admin").First(&adminRole).Error; err == nil {
		bindRoleResources(db, adminRole, all, func(model.Resource) bool { return true })
	}

	// Add query resource access permissions to guest
	var guestRole model.Role
	if err := db.Where("name", "guest").First(&guestRole).Error; err == nil {
		bindRoleResources(db, guestRole, all, func(r model.Resource) bool { return r.Method == "GET" })
	}

	slog.Info("-----Initialize interface resources end-----")
}

// 把满足 match 的资源挂到角色下, 已经存在的关联跳过
func bindRoleResources(db *gorm.DB, role model.Role, resources []model.Resource, match func(model.Resource) bool) {
	for _, resource := range resources {
		if resource.ID == 0 || !match(resource) {
			continue
		}
		err := db.Create(&model.RoleResource{RoleId: role.ID, ResourceId: resource.ID}).Error
		if err != nil && !isDuplicate(err) {
			slog.Error(role.Name + " 角色资源关联关系初始化失败" + err.Error())
		}
	}
}

// Generate default menus
func generateDefaultMenus(db *gorm.DB) {
	slog.Info("-----Initialize menus start-----")

	parents := []model.Menu{
		{Name: "Trang chủ", Path: "/home", Icon: "ic:sharp-home", OrderNum: 0, Component: "/home", Redirect: "/home", Catalogue: true}, // catalogue
		{Name: "Quản lý bài viết", Path: "/article", Icon: "ic:twotone-article", OrderNum: 1, Component: "Layout", Redirect: "/article/list"},
		{Name: "Quản lý phân quyền", Path: "/auth", Icon: "cib:adguard", OrderNum: 3, Component: "Layout", Redirect: "/auth/menu"},
		{Name: "Quản lý tin nhắn", Path: "/message", Icon: "ic:twotone-email", OrderNum: 2, Component: "Layout", Redirect: "/message/comment"},
		{Name: "Quản lý người dùng", Path: "/user", Icon: "ph:user-list-bold", OrderNum: 4, Component: "Layout", Redirect: "/user/list"},
		{Name: "Quản lý log", Path: "/log", Icon: "material-symbols:receipt-long-outline-rounded", OrderNum: 6, Component: "Layout", Redirect: "/log/operation"},
		{Name: "Quản lý hệ thống", Path: "/setting", Icon: "ion:md-settings", OrderNum: 5, Component: "Layout", Redirect: "/setting/website"},
		{Name: "Trung tâm cá nhân", Path: "/profile", Icon: "mdi:account", OrderNum: 7, Component: "/profile", Redirect: "/profile", Catalogue: true}, // catalogue
	}

	for i := range parents {
		if err := upsertSeedParentMenu(db, &parents[i]); err != nil {
			slog.Error(parents[i].Name + " menu initialization failed: " + err.Error())
		}
	}

	menus := []model.Menu{
		{Name: "Đăng bài viết", Path: "write", Component: "/article/write", Icon: "icon-park-outline:write", OrderNum: 1, ParentId: parents[1].ID},
		{Name: "Danh sách bài viết", Path: "list", Component: "/article/list", Icon: "material-symbols:format-list-bulleted", OrderNum: 2, ParentId: parents[1].ID},
		{Name: "Quản lý danh mục", Path: "category", Component: "/article/category", Icon: "tabler:category", OrderNum: 3, ParentId: parents[1].ID},
		{Name: "Quản lý tag", Path: "tag", Component: "/article/tag", Icon: "tabler:tag", OrderNum: 4, ParentId: parents[1].ID},
		{Name: "Quản lý bài đăng ngắn", Path: "talk", Component: "/article/talk", Icon: "mdi:message-text-outline", OrderNum: 5, ParentId: parents[1].ID},
		{Name: "Sửa bài viết", Path: "write/:id", Component: "/article/write", Icon: "icon-park-outline:write", OrderNum: 1, ParentId: parents[1].ID, Hidden: true},
		{Name: "Quản lý menu", Path: "menu", Component: "/auth/menu", Icon: "ic:twotone-menu-book", OrderNum: 1, ParentId: parents[2].ID},
		{Name: "Quản lý API", Path: "resource", Component: "/auth/resource", Icon: "mdi:api", OrderNum: 2, ParentId: parents[2].ID},
		{Name: "Quản lý role", Path: "role", Component: "/auth/role", Icon: "carbon:user-role", OrderNum: 3, ParentId: parents[2].ID},
		{Name: "Quản lý bình luận", Path: "comment", Component: "/message/comment", Icon: "ic:twotone-comment", OrderNum: 1, ParentId: parents[3].ID},
		{Name: "Quản lý lời nhắn", Path: "leave-msg", Component: "/message/leave-msg", Icon: "ic:twotone-message", OrderNum: 2, ParentId: parents[3].ID},
		{Name: "Danh sách người dùng", Path: "list", Component: "/user/list", Icon: "mdi:account", OrderNum: 1, ParentId: parents[4].ID},
		{Name: "Người dùng online", Path: "online", Component: "/user/online", Icon: "ic:outline-online-prediction", OrderNum: 2, ParentId: parents[4].ID},
		{Name: "Operation log", Path: "operation", Component: "/log/operation", Icon: "mdi:book-open-page-variant-outline", OrderNum: 1, ParentId: parents[5].ID},
		{Name: "Login log", Path: "login", Component: "/log/login", Icon: "material-symbols:login", OrderNum: 2, ParentId: parents[5].ID},
		{Name: "Lỗi frontend", Path: "error", Component: "/log/error", Icon: "mdi:bug-outline", OrderNum: 3, ParentId: parents[5].ID},
		{Name: "Cấu hình website", Path: "website", Component: "/setting/website", Icon: "el:website", OrderNum: 1, ParentId: parents[6].ID},
		{Name: "Quản lý trang", Path: "page", Component: "/setting/page", Icon: "iconoir:journal-page", OrderNum: 2, ParentId: parents[6].ID},
		{Name: "Quản lý friend link", Path: "link", Component: "/setting/link", Icon: "mdi:telegram", OrderNum: 3, ParentId: parents[6].ID},
		{Name: "Giới thiệu", Path: "about", Component: "/setting/about", Icon: "cib:about-me", OrderNum: 4, ParentId: parents[6].ID},
	}

	for i := range menus {
		if err := upsertSeedChildMenu(db, &menus[i]); err != nil {
			slog.Error(menus[i].Name + " menu initialization failed: " + err.Error())
		}
	}

	// 给 admin 和 guest 角色添加所有菜单访问权限
	for _, name := range []string{"admin", "guest"} {
		var role model.Role
		if err := db.Where("name", name).First(&role).Error; err != nil {
			continue
		}
		bindRoleMenus(db, role, menus)
	}

	slog.Info("-----Initialize menus end-----")
}

// Menu names are display text and may change when translating the UI. Seed by
// stable route path instead of name, otherwise every translation creates a new
// route tree beside the old one.
func upsertSeedParentMenu(db *gorm.DB, desired *model.Menu) error {
	var matches []model.Menu
	if err := db.Where("parent_id = 0 AND path = ?", desired.Path).Order("id").Find(&matches).Error; err != nil {
		return err
	}
	if len(matches) == 0 {
		return db.Create(desired).Error
	}

	canonical := matches[0]
	for _, duplicate := range matches[1:] {
		if err := db.Model(&model.Menu{}).Where("parent_id = ?", duplicate.ID).Update("parent_id", canonical.ID).Error; err != nil {
			return err
		}
		if err := mergeMenuRoleBindings(db, canonical.ID, duplicate.ID); err != nil {
			return err
		}
		if err := db.Delete(&duplicate).Error; err != nil {
			return err
		}
	}

	desired.ID = canonical.ID
	return model.SaveOrUpdateMenu(db, desired)
}

func upsertSeedChildMenu(db *gorm.DB, desired *model.Menu) error {
	var matches []model.Menu
	// parent_id = 0 also catches orphan rows produced when an older parent seed
	// failed (for example /menu after "Permission Management" exceeded varchar(20)).
	if err := db.Where("path = ? AND (parent_id = ? OR (parent_id = 0 AND component = ?))", desired.Path, desired.ParentId, desired.Component).
		Order("id").Find(&matches).Error; err != nil {
		return err
	}
	if len(matches) == 0 {
		return db.Create(desired).Error
	}

	canonical := matches[0]
	for _, duplicate := range matches[1:] {
		if err := mergeMenuRoleBindings(db, canonical.ID, duplicate.ID); err != nil {
			return err
		}
		if err := db.Delete(&duplicate).Error; err != nil {
			return err
		}
	}

	desired.ID = canonical.ID
	return model.SaveOrUpdateMenu(db, desired)
}

func mergeMenuRoleBindings(db *gorm.DB, keepMenuID, removeMenuID int) error {
	var roleIDs []int
	if err := db.Model(&model.RoleMenu{}).Where("menu_id = ?", removeMenuID).Pluck("role_id", &roleIDs).Error; err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		binding := model.RoleMenu{RoleId: roleID, MenuId: keepMenuID}
		if err := db.Where(binding).FirstOrCreate(&binding).Error; err != nil {
			return err
		}
	}
	return db.Where("menu_id = ?", removeMenuID).Delete(&model.RoleMenu{}).Error
}

// 把菜单挂到角色下, 已经存在的关联跳过
// 重复执行是正常的(每次容器启动都会跑一遍), 不要为此刷一堆日志
func bindRoleMenus(db *gorm.DB, role model.Role, menus []model.Menu) {
	for _, menu := range menus {
		if menu.ID == 0 {
			continue
		}
		err := db.Create(&model.RoleMenu{RoleId: role.ID, MenuId: menu.ID}).Error
		if err != nil && !isDuplicate(err) {
			slog.Error(role.Name + " 角色菜单关联关系初始化失败" + err.Error())
		}
	}
}

// sqlite 和 MySQL 的唯一约束冲突错误文案不同, 统一判断
// 种子数据允许重复执行, 冲突说明已经初始化过, 不算失败
func isDuplicate(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "Duplicate entry")
}
