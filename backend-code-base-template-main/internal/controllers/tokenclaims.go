package controllers

import (
	"backend-code-base-template/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenClaimsController struct {
	router *gin.RouterGroup
}

func NewTokenClaimsController(router *gin.RouterGroup) *TokenClaimsController {
	return &TokenClaimsController{
		router: router,
	}
}

func (oauth *TokenClaimsController) InitRoutes() {
	oauth.router.GET("/sso/claims", func(ctx *gin.Context) {
		token := ctx.Query("token")
		respTokenClaims := utilities.ValidateJwtToken(token)
		ctx.JSON(http.StatusOK, respTokenClaims)
	})

}
