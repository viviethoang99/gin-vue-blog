package handle

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type BlogInfo struct{}

type BlogHomeVO struct {
	ArticleCount int `json:"article_count"` // Article count
	UserCount    int `json:"user_count"`    // User count
	MessageCount int `json:"message_count"` // Message count
	ViewCount    int `json:"view_count"`    // Visit count
	// CategoryCount int64 `json:"category_count"` // Category count
	// TagCount      int64 `json:"tag_count"`      // Tag count
	// BlogConfig    model.BlogConfigDetail `json:"blog_config"`    // Blog info
	// PageList      []Page                 `json:"pageList"`
}

type AboutReq struct {
	Content string `json:"content"`
}

func (*BlogInfo) GetConfigMap(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	// get from redis cache
	cache, err := getConfigCache(rdb)
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	if len(cache) > 0 {
		slog.Debug("get config from redis cache")
		ReturnSuccess(c, cache)
		return
	}

	// get from db
	data, err := model.GetConfigMap(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// add to redis cache
	if err := addConfigCache(rdb, data); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, data)
}

func (*BlogInfo) UpdateConfig(c *gin.Context) {
	var m map[string]string
	if err := c.ShouldBindJSON(&m); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := model.CheckConfigMap(GetDB(c), m); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// delete cache
	if err := removeConfigCache(GetRDB(c)); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// @Summary Get blog homepage info
// @Description Get blog homepage info
// @Tags blog_info
// @Produce json
// @Success 0 {object} Response[model.BlogHomeVO]
// @Router /home [get]
func (*BlogInfo) GetHomeInfo(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	articleCount, err := model.Count(db, &model.Article{}, "status = ? AND is_delete = ?", 1, 0)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	userCount, err := model.Count(db, &model.UserInfo{})
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	messageCount, err := model.Count(db, &model.Message{})
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	viewCount, err := rdb.Get(rctx, g.VIEW_COUNT).Int()
	if err != nil && err != redis.Nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, BlogHomeVO{
		ArticleCount: articleCount,
		UserCount:    userCount,
		MessageCount: messageCount,
		ViewCount:    viewCount,
	})
}

// @Summary Get About
// @Description Get About
// @Tags blog_info
// @Produce json
// @Success 0 {object} Response[string]
// @Router /about [get]
func (*BlogInfo) GetAbout(c *gin.Context) {
	ReturnSuccess(c, model.GetConfig(GetDB(c), g.CONFIG_ABOUT))
}

// @Summary Update About
// @Description Update About
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "About"
// @Success 0 {object} Response[string]
// @Router /about [put]
func (*BlogInfo) UpdateAbout(c *gin.Context) {
	var req AboutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.CheckConfig(GetDB(c), g.CONFIG_ABOUT, req.Content)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, req.Content)
}

// @Summary Report user info
// @Description Report when user logs into admin
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "User info"
// @Success 0 {object} Response[any]
// @Router /report [post]
func (*BlogInfo) Report(c *gin.Context) {
	rdb := GetRDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	userAgent := utils.IP.GetUserAgent(c)
	browser := userAgent.Name + " " + userAgent.Version.String()
	os := userAgent.OS + " " + userAgent.OSVersion.String()
	uuid := utils.MD5(ipAddress + browser + os)

	ctx := context.Background()

	// Current user not counted in visitor set
	if !rdb.SIsMember(ctx, g.KEY_UNIQUE_VISITOR_SET, uuid).Val() {
		// Collect area statistics
		ipSource := utils.IP.GetIpSource(ipAddress)
		if ipSource != "" { // Got a specific location, extract province
			address := strings.Split(ipSource, "|")
			province := strings.ReplaceAll(address[2], "省", "")
			rdb.HIncrBy(ctx, g.VISITOR_AREA, province, 1)
		} else {
			rdb.HIncrBy(ctx, g.VISITOR_AREA, "Unknown", 1)
		}
		// Increase visit count by 1
		rdb.Incr(ctx, g.VIEW_COUNT)
		// Add current user to visitor set
		rdb.SAdd(ctx, g.KEY_UNIQUE_VISITOR_SET, uuid)
	}

	ReturnSuccess(c, nil)
}

// Get blog settings
// func GetBlogConfig() model.BlogConfigDetail {
// 	// Try to get value from Redis
// 	blogConfig := utils.Redis.GetVal(KEY_BLOG_CONFIG)
// 	// If Redis missing, query DB and set Redis
// 	if blogConfig == "" {
// 		blogConfig = dao.GetOne(model.BlogConfig{}, "id", 1).Config
// 		utils.Redis.Set(KEY_BLOG_CONFIG, blogConfig, 0)
// 	}
// 	// Deserialize string to Go object
// 	var result model.BlogConfigDetail
// 	utils.Json.Unmarshal(blogConfig, &result)
// 	return result
// }
