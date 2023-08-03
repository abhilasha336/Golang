package utilities

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetVersionFromContext
// get the accepted version from header, this should be implemented as middleware
func GetVersionFromContext(ctx *gin.Context) string {
	return ctx.GetHeader("Accept-version")
}

// PrepareVersionName
// this function will prepare a version name v1.1 into v1_1
func PrepareVersionName(version string) string {
	version = strings.ReplaceAll(version, ".", "_")
	return version
}

// get Context
// read the context data and type assert into curresponding concrete value
func GetContext[T any](ctx *gin.Context, name string) (T, bool) {
	value, exists := ctx.Get(name)
	return value.(T), exists
}
