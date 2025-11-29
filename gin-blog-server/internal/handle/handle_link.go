package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Link struct{}

// Add or edit friend link
type AddOrEditLinkReq struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	Avatar  string `json:"avatar"`
	Address string `json:"address" binding:"required"`
	Intro   string `json:"intro"`
}

// @Summary Get friend link list
// @Description Get friend link list by conditions
// @Tags Link
// @Param page_size query int false "Current page"
// @Param page_num query int false "Page size"
// @Param keyword query string false "Keyword"
// @Accept json
// @Produce json
// @Success 0 {object} Response[PageResult[model.FriendLink]]
// @Security ApiKeyAuth
// @Router /link/list [get]
func (*Link) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	data, total, err := model.GetLinkList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.FriendLink]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// @Summary Add or edit friend link
// @Description Add or edit friend link
// @Tags Link
// @Param form body AddOrEditLinkReq true "Add or edit friend link"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.FriendLink]
// @Security ApiKeyAuth
// @Router /link [post]
func (*Link) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	link, err := model.SaveOrUpdateLink(GetDB(c), req.ID, req.Name, req.Avatar, req.Address, req.Intro)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, link)
}

// @Summary Delete friend links (batch)
// @Description Delete friend links by ID array
// @Tags Link
// @Param ids body []int true "Friend link ID array"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int64]
// @Security ApiKeyAuth
// @Router /link [delete]
func (*Link) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(&model.FriendLink{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, g.ErrDbOp, result.Error)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
