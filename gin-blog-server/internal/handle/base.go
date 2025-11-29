package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

/*
Response design:
- Do not use HTTP status codes to represent business status; use business codes instead
- Any request that reaches the backend returns HTTP 200
- Business code 0 indicates success; non-zero indicates failure
- Only when a backend panic occurs and is caught by gin middleware, HTTP 500 is returned
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
	Page    int    `form:"page_num"`  // Current page (starts from 1)
	Size    int    `form:"page_size"` // Page size
	Keyword string `form:"keyword"`   // Keyword
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
