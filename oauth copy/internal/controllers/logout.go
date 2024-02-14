package controllers

import (
	"context"
	"net/http"
	"oauth/internal/entities"
	"oauth/utilities"

	"github.com/gin-gonic/gin"
	log "gitlab.com/tuneverse/toolkit/core/logger"
)

// fn helps to revoke token
func (oauth *OauthController) OauthLogOut(ctx *gin.Context) {

	var (
		refreshToken entities.Refresh
		log          = log.Log().WithContext(ctx)
	)

	cfg := oauth.cfg
	accessToken := ctx.GetHeader("Authorization")
	partnerID := ctx.GetHeader("partner_id")

	if err := ctx.ShouldBindJSON(&refreshToken); err != nil {
		log.Errorf("OauthLogOut controller-payload bind failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"response": "error in payload"})
		return
	}

	responseToken := utilities.ValidateJwtToken(accessToken, cfg.JwtKey)
	err := oauth.useCase.Logout(context.Background(), refreshToken, accessToken, partnerID, *responseToken.MemberID)
	if err != nil {
		log.Errorf("OauthLogOut controller-unable to revoke token status in db: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"response": "logout failed"})
		log.Printf("log out failed")
		return
	}

	response := entities.Response{
		Error:   nil,
		Message: "logged out",
		Data: []map[string]interface{}{
			{
				"message": "loggedout successfully",
			},
		},
	}

	ctx.JSON(http.StatusOK, response)
	log.Printf("logout successfull")

}
