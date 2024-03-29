package utils

import (
	"github.com/gin-gonic/gin"
)

// get Context
// read the context data and type assert into curresponding concrete value
func GetContext[T any](ctx *gin.Context, name string) (T, bool) {
	value, exists := ctx.Get(name)
	return value.(T), exists
}

// Function to retrieve header data
func GetHeader(ctx *gin.Context, header string) string {
	value := ctx.GetHeader(header)
	return value
}
