package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

/*
响应设计方案：不使用 HTTP 码来表示业务状态, 采用业务状态码的方式
- 只要能到达后端的请求, HTTP 状态码都为 200
- 业务状态码为 0 表示成功, 其他都表示失败
- 当后端发生 panic 并且被 gin 中间件捕获时, 才会返回 HTTP 500 状态码
*/

// Response structure
type Response[T any] struct {
	Code    int    `json:"code"`    // Business status code
	Message string `json:"message"` // Response message
	Data    T      `json:"data"`    // Response data
}

// HTTP code + Business code + Message + Data
func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// Business code + Data
func ReturnResponse(c *gin.Context, r g.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

// Success business code + Data
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, g.OkResult, data)
}

// All predictable errors = business + system errors, handled at business layer, return HTTP 200
// Unpredictable errors trigger panic, caught by gin middleware, return HTTP 500
// err is business error; data is error payload (error or string)
func ReturnError(c *gin.Context, r g.Result, data any) {
	slog.Info("[Func-ReturnError] " + r.Msg())

	var val string = r.Msg()

	if data != nil {
		switch v := data.(type) {
		case error:
			val = v.Error()
		case string:
			val = v
		}
		slog.Error(val) // Error log
	}

	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code:    r.Code(),
			Message: r.Msg(),
			Data:    val,
		},
	)
}

// Pagination query parameters
type PageQuery struct {
	Page int `form:"page_num"` // 当前页数（从1开始）
	// 每页条数; 限制成正数, 避免请求参数传进 model.PageSizeAll 变成全表查询
	Size    int    `form:"page_size" binding:"omitempty,min=1"`
	Keyword string `form:"keyword"` // 搜索关键字
}

// Pagination response data
type PageResult[T any] struct {
	Page  int   `json:"page_num"`  // Current page
	Size  int   `json:"page_size"` // Page size
	Total int64 `json:"total"`     // Total count
	List  []T   `json:"page_data"` // Paged data
}

// Get *gorm.DB
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}

// Get *redis.Client
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CTX_RDB).(*redis.Client)
}

/*
获取当前登录用户, 未登录时直接响应错误并返回 false

调用方必须在返回 false 时立即 return。
JWTAuth 中间件对资源表中不存在的接口会跳过鉴权, 因此 handler 里
不能假设一定拿得到用户, 否则会 nil 解引用导致 panic。
*/
func MustCurrentUserAuth(c *gin.Context) (*model.UserAuth, bool) {
	auth, err := CurrentUserAuth(c)
	if err != nil || auth == nil {
		ReturnError(c, g.ErrTokenNotExist, err)
		return nil, false
	}
	return auth, true
}

/*
Get current logged-in user info
1. If a user object exists on gin Context, it was already fetched in this request chain
2. Get uid from session
3. Fetch user info by uid and set it on gin Context
*/
func CurrentUserAuth(c *gin.Context) (*model.UserAuth, error) {
	key := g.CTX_USER_AUTH

	// 1
	if cache, exist := c.Get(key); exist && cache != nil {
		slog.Debug("[Func-CurrentUserAuth] get from cache: " + cache.(*model.UserAuth).Username)
		return cache.(*model.UserAuth), nil
	}

	// 2
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("no user_auth_id in session")
	}

	// 3
	db := GetDB(c)
	user, err := model.GetUserAuthInfoById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}

/*
访客指纹: IP + 浏览器 + 操作系统 的哈希, 用于同一访客的去重统计
*/
func visitorFingerprint(c *gin.Context) string {
	// 非浏览器(如 curl)或没有 User-Agent 时, 解析结果为 nil, 不能直接取字段
	var browser, os string
	if ua := utils.IP.GetUserAgent(c); ua != nil {
		browser = ua.Name + " " + ua.Version.String()
		os = ua.OS + " " + ua.OSVersion.String()
	}

	return utils.MD5(clientIP(c) + browser + os)
}

/*
客户端 IP, 已剥掉端口

RemoteAddr 形如 "1.2.3.4:54321", 端口每次请求都变, 不剥掉的话按 IP 做的
去重和限流都会失效(同一个来源被当成无数个新来源)。
*/
func clientIP(c *gin.Context) string {
	ipAddress := utils.IP.GetIpAddress(c)
	if host, _, err := net.SplitHostPort(ipAddress); err == nil {
		return host
	}
	return ipAddress
}
