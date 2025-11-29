package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type Page struct{}

// @Summary Get page list
// @Description Get page list by conditions
// @Tags Page
// @Accept json
// @Produce json
// @Success 0 {object} Response[[]model.Page]
// @Security ApiKeyAuth
// @Router /page/list [get]
func (*Page) GetList(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	// get from cache
	cache, err := getPageCache(rdb)
	if cache != nil && err == nil {
		slog.Debug("[handle-page-GetList] get page list from cache")
		ReturnSuccess(c, cache)
		return
	}

	switch err {
	case redis.Nil:
		break
	default:
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	// get from db
	data, _, err := model.GetPageList(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// add to cache
	if err := addPageCache(GetRDB(c), data); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, data)
}

// @Summary Add or edit page
// @Description Add or edit page
// @Tags Page
// @Param form body AddOrEditPageReq true "Add or edit page"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.Page]
// @Security ApiKeyAuth
// @Router /page [post]
func (*Page) SaveOrUpdate(c *gin.Context) {
	var req model.Page
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	page, err := model.SaveOrUpdatePage(db, req.ID, req.Name, req.Label, req.Cover)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// delete cache
	if err := removePageCache(rdb); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, page)
}

// @Summary Delete pages
// @Description Delete pages by ID array
// @Tags Page
// @Param ids body []int true "Page ID array"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /page [delete]
func (*Page) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(model.Page{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	// delete cache
	if err := removePageCache(GetRDB(c)); err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
