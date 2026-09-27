package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"log/slog"

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
// @Accept json
// @Produce json
// @Param ids body []int true "Comment ID array"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /comment [delete]
func (*Comment) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, deletedIds, err := model.DeleteComments(GetDB(c), ids)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 与文章删除保持一致: 清掉点赞计数, 否则新评论复用 id 会继承旧数据
	cleanCommentCounters(GetRDB(c), deletedIds)

	ReturnSuccess(c, rows)
}

// @Summary Update comment review (batch)
// @Description Update review status by ID array
// @Tags Comment
// @Accept json
// @Produce json
// @Param form body UpdateReviewReq true "Update review status"
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /comment/review [put]
func (*Comment) UpdateReview(c *gin.Context) {
	var req UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	db := GetDB(c)
	rows, approved, err := model.ReviewComments(db, req.Ids, req.IsReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 过审这一刻才补发通知: 评论待审核期间是不可见的, 当时没发
	// 通知写失败不影响审核结果, 只记日志
	for i := range approved {
		if err := model.NotifyOnComment(db, &approved[i]); err != nil {
			slog.Warn("审核通过后写站内通知失败", "err", err, "comment_id", approved[i].ID)
		}
	}

	ReturnSuccess(c, rows)
}

// @Summary Query comment list
// @Description 支持按昵称/审核状态/类型过滤
// @Tags Comment
// @Produce json
// @Param nickname query string false "Nickname"
// @Param is_review query bool false "Review status"
// @Param type query int false "评论类型(1-文章 2-友链 3-说说)"
// @Param page_num query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 0 {object} Response[PageResult[model.Comment]]
// @Security ApiKeyAuth
// @Router /comment/list [get]
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
