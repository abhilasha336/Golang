package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler function used too check health status
func (oauth *OauthController) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "failure",
		"message": "failure",
	})
}
