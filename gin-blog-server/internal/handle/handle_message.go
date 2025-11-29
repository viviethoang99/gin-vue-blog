package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

// Update review status (batch)
type UpdateReviewReq struct {
	Ids      []int `json:"ids"`
	IsReview bool  `json:"is_review"`
}

// Query list with conditions
type MessageQuery struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *bool  `form:"is_review"`
}

type Message struct{}

// @Summary Delete messages (batch)
// @Description Delete messages by ID array
// @Tags Category
// @Param ids body []int true "Message ID array"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /category [delete]
func (*Message) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.DeleteMessages(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// @Summary Update message review (batch)
// @Description Update review status by ID array
// @Tags Message
// @Param form body UpdateReviewReq true "Update review status"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /message/review [put]
func (*Message) UpdateReview(c *gin.Context) {
	var req UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.UpdateMessagesReview(GetDB(c), req.Ids, req.IsReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// @Summary Query message list
// @Description Get message list by conditions
// @Tags Message
// @Param nickname query string false "Nickname"
// @Param is_review query int false "Review status"
// @Param page_size query int false "Current page"
// @Param page_num query int false "Page size"
// @Accept json
// @Produce json
// @Success 0 {object} Response[PageResult[model.Message]]
// @Security ApiKeyAuth
// @Router /message/list [get]
func (*Message) GetList(c *gin.Context) {
	var query MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetMessageList(GetDB(c), query.Page, query.Size, query.Nickname, query.IsReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, PageResult[model.Message]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}
