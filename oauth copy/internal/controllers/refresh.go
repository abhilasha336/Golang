package controllers

import (
	"context"
	"net/http"
	"oauth/internal/consts"
	"oauth/internal/entities"
	"oauth/utilities"

	"github.com/gin-gonic/gin"
	log "gitlab.com/tuneverse/toolkit/core/logger"
)

// function which generates new token and invalidate old jwt refresh token
func (oauth *OauthController) OauthRefresh(ctx *gin.Context) {

	var (
		log = log.Log().WithContext(ctx)
	)
	cfg := oauth.cfg

	partnerID := ctx.GetHeader("partner_id")
	oldRefresh := ctx.GetHeader("Authorization")

	responseData := utilities.ValidateJwtToken(oldRefresh, cfg.JwtKey)

	jwtPayload := entities.OAuthData{}
	jwtPayload.MemberEmail = responseData.MemberEmail
	jwtPayload.MemberID = responseData.MemberID
	jwtPayload.PartnerID = responseData.PartnerID
	jwtPayload.MemberName = responseData.MemberName
	jwtPayload.MemberType = responseData.MemberType
	jwtPayload.Roles = responseData.Roles
	jwtPayload.PartnerName = responseData.PartnerName

	newToken := utilities.GenerateJwtToken(jwtPayload, consts.ExpTime, cfg.JwtKey)
	newRefreshTok := utilities.GenerateJwtToken(jwtPayload, consts.RefExpTime, cfg.JwtKey)

	err := oauth.useCase.DeleteAndInsertRefreshToken(context.Background(), oldRefresh, newToken, newRefreshTok, partnerID, jwtPayload.MemberID)
	if err != nil {
		log.Errorf("OauthRefresh controller-token deletion and updation failed: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"response": "unauthorised,invalid token"})
		return
	}

	response := entities.Response{
		Error:   err,
		Message: "token regenerated successfully",
		Data: []map[string]interface{}{
			{
				"token":   newToken,
				"refresh": newRefreshTok,
			},
		},
	}
	log.Printf("oauth new refreshtoken generation successful")
	ctx.JSON(http.StatusOK, response)
}
