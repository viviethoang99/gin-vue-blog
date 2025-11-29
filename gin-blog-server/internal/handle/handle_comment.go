package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type CommentQuery struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *bool  `form:"is_review"`
	Type     int    `form:"type"`
}

type Comment struct{}

// @Summary Delete comments (batch)
// @Description Delete comments by ID array
// @Tags Comment
// @Param ids body []int true "Comment ID array"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /comment [delete]
func (*Comment) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(model.Comment{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}

// @Summary Update comment review (batch)
// @Description Update review status by ID array
// @Tags Comment
// @Param form body UpdateReviewReq true "Update review status"
// @Accept json
// @Produce json
// @Success 0 {object} Response[any]
// @Security ApiKeyAuth
// @Router /comment/review [put]
func (*Comment) UpdateReview(c *gin.Context) {
	var req UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	maps := map[string]any{"is_review": req.IsReview}
	result := GetDB(c).Model(model.Comment{}).Where("id in ?", req.Ids).Updates(maps)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}

// @Summary Query comment list
// @Description Get comment list by conditions
// @Tags Comment
// @Param nickname query string false "Nickname"
// @Param is_review query int false "Review status"
// @Param type query int false "Comment type"
// @Param page_num query int false "Page number"
// @Param page_size query int false "Page size"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.CommentVO]
// @Security ApiKeyAuth
// @Router /comment [get]
func (*Comment) GetList(c *gin.Context) {
	var query CommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	list, total, err := model.GetCommentList(GetDB(c), query.Page, query.Size, query.Type, query.IsReview, query.Nickname)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.Comment]{
		Total: total,
		List:  list,
		Size:  query.Size,
		Page:  query.Page,
	})

}
