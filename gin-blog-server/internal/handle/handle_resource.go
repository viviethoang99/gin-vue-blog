package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Resource struct{}

type TreeOptionVO struct {
	ID       int            `json:"key"`
	Label    string         `json:"label"`
	Children []TreeOptionVO `json:"children"`
}

type ResourceTreeVO struct {
	ID        int              `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Name      string           `json:"name"`
	Url       string           `json:"url"`
	Method    string           `json:"request_method"`
	Anonymous bool             `json:"is_anonymous"`
	Children  []ResourceTreeVO `json:"children"`
}

// TODO: Use oneof tag to validate data
type AddOrEditResourceReq struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Method   string `json:"request_method"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
}

type EditAnonymousReq struct {
	ID        int  `json:"id" binding:"required"`
	Anonymous bool `json:"is_anonymous"`
}

// Get resource list (tree)
func (*Resource) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")

	resourceList, err := model.GetResourceList(GetDB(c), keyword)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, resources2ResourceVos(resourceList))
}

// Get option data (tree)
func (*Resource) GetOption(c *gin.Context) {
	result := make([]TreeOptionVO, 0)

	db := GetDB(c)
	resources, err := model.GetResourceList(db, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)

	for _, item := range parentList {
		var children []TreeOptionVO
		for _, re := range childrenMap[item.ID] {
			children = append(children, TreeOptionVO{
				ID:    re.ID,
				Label: re.Name,
			})
		}
		result = append(result, TreeOptionVO{
			ID:       item.ID,
			Label:    item.Name,
			Children: children,
		})
	}
	ReturnSuccess(c, result)
}

// Add or edit resource; update casbin_rule accordingly
func (*Resource) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditResourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	err := model.SaveOrUpdateResource(db, req.ID, req.ParentId, req.Name, req.Url, req.Method)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// Edit resource anonymous access; update casbin_rule accordingly
func (*Resource) UpdateAnonymous(c *gin.Context) {
	var req EditAnonymousReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.UpdateResourceAnonymous(GetDB(c), req.ID, req.Anonymous)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// TODO: Consider behavior after deleting a module; currently cannot delete if it has children
// TODO: Force delete?
func (*Resource) Delete(c *gin.Context) {
	resourceId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// Check if the resource is used by any role
	use, _ := model.CheckResourceInUse(db, resourceId)
	if use {
		ReturnError(c, g.ErrResourceUsedByRole, nil)
		return
	}

	// Get the resource
	resource, err := model.GetResourceById(db, resourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrResourceNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// If it's a module, check whether it has child resources
	if resource.ParentId == 0 {
		hasChild, _ := model.CheckResourceHasChild(db, resourceId)
		if hasChild {
			ReturnError(c, g.ErrResourceHasChildren, nil)
			return
		}
	}

	rows, err := model.DeleteResource(db, resourceId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// []Resource => []ResourceVO
func resources2ResourceVos(resources []model.Resource) []ResourceTreeVO {
	list := make([]ResourceTreeVO, 0)
	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)
	for _, item := range parentList {
		resourceVO := resource2ResourceVo(item)
		resourceVO.Children = make([]ResourceTreeVO, 0)
		for _, child := range childrenMap[item.ID] {
			resourceVO.Children = append(resourceVO.Children, resource2ResourceVo(child))
		}
		list = append(list, resourceVO)
	}
	return list
}

// Resource => ResourceVO
func resource2ResourceVo(r model.Resource) ResourceTreeVO {
	return ResourceTreeVO{
		ID:        r.ID,
		Name:      r.Name,
		Url:       r.Url,
		Method:    r.Method,
		Anonymous: r.Anonymous,
		CreatedAt: r.CreatedAt,
	}
}

// Map storing child resource list for each node
// key: resourceId
// value: childrenList
func getChildrenMap(resources []model.Resource) map[int][]model.Resource {
	m := make(map[int][]model.Resource)
	for _, r := range resources {
		if r.ParentId != 0 {
			m[r.ParentId] = append(m[r.ParentId], r)
		}
	}
	return m
}

// Get top-level resources (parent_id == 0)
func getModuleList(resources []model.Resource) []model.Resource {
	list := make([]model.Resource, 0)
	for _, r := range resources {
		if r.ParentId == 0 {
			list = append(list, r)
		}
	}
	return list
}
