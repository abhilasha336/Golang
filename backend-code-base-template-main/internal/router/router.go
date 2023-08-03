package router

import "github.com/gin-gonic/gin"

type RouterGroup struct {
	groups []*gin.RouterGroup
}

func NewGroupGroup(groups []*gin.RouterGroup) RouterGroup {
	return RouterGroup{
		groups,
	}
}

func (g *RouterGroup) Handle(method string, path string, handler gin.HandlerFunc) {
	for _, group := range g.groups {
		group.Handle(method, path, handler)
	}
}

func InitVersions(router *gin.Engine) RouterGroup {
	v1 := router.Group("/v1")
	v2 := router.Group("/v2")

	return NewGroupGroup([]*gin.RouterGroup{v1, v2})
}
