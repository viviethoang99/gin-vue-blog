package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

type AddOrEditTagReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

// @Summary Get tag list
// @Description Get tag list by conditions
// @Tags Tag
// @Param page_size query int false "Current page"
// @Param page_num query int false "Page size"
// @Param keyword query string false "Keyword"
// @Accept json
// @Produce json
// @Success 0 {object} Response[PageResult[model.TagVO]] "Success"
// @Security ApiKeyAuth
// @Router /tag/list [get]
func (*Tag) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetTagList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.TagVO]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary Add or edit tag
// @Description Add or edit tag
// @Tags Tag
// @Param form body AddOrEditTagReq true "Add or edit tag"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /tag [post]
func (*Tag) SaveOrUpdate(c *gin.Context) {
	var form AddOrEditTagReq
	if err := c.ShouldBindJSON(&form); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	tag, err := model.SaveOrUpdateTag(GetDB(c), form.ID, form.Name)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, tag)
}

// TODO: Delete behavior: add force delete to remove related data if exists
// @Summary Delete tags (batch)
// @Description Delete tags by ID array
// @Tags Tag
// @Param ids body []int true "Tag ID array"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /tag [delete]
func (*Tag) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	db := GetDB(c)

	// Check if the tag has related articles
	count, err := model.Count(db, &model.ArticleTag{}, "tag_id in ?", ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	if count > 0 {
		// ReturnError(c, g.ERROR_TAG_ART_EXIST, nil)
		ReturnError(c, g.ErrTagHasArt, nil)
		return
	}

	result := db.Delete(model.Tag{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}

// @Summary Get tag options
// @Description Get tag options list
// @Tags Tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /tag/option [get]
// @Success 0 {object} Response[model.OptionVO]
func (*Tag) GetOption(c *gin.Context) {
	list, err := model.GetTagOption(GetDB(c))
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
