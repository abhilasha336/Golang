package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"oauth/internal/consts"
	"oauth/internal/entities"
	"oauth/utilities"

	log "gitlab.com/tuneverse/toolkit/core/logger"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/utils"
)

// controller function which handles manual login with email and password(login as internal server)
func (oauth *OauthController) OauthLogIn(ctx *gin.Context) {

	var (
		loginRequest     entities.LoginRequest
		apiResponseLogin entities.BasicMemberDataResponse
		refTok           entities.Refresh
		log              = log.Log().WithContext(ctx)
		err              error
		tokenPayload     entities.OAuthData
		refToken         entities.Refresh
	)

	cfg := oauth.cfg

	provider := ctx.GetHeader("provider")
	partnerID := ctx.GetHeader("partner_id")

	// Bind the request body to the LoginRequest struct
	if err := ctx.ShouldBind(&loginRequest); err != nil {
		log.Errorf("OauthLogIn contoller-bind request-payload error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please check payload"})
		return
	}

	//member get api call request body
	body := map[string]interface{}{
		"email":    loginRequest.Email,
		"provider": provider,
		"password": loginRequest.Password,
	}

	tokenPayload.PartnerID = partnerID
	tokenPayload.Email = loginRequest.Email
	getMemberToken := utilities.GenerateJwtToken(tokenPayload, consts.TempExpTime, cfg.JwtKey)
	refToken.RefreshToken = getMemberToken

	err = oauth.useCase.PostRefreshToken(context.Background(), refToken, getMemberToken, partnerID, nil)
	if err != nil {
		log.Errorf("OauthLogIn contoller-dummy token entry in refreshtoken table failed: %v", err)
	}

	//header to call member api
	header := map[string]interface{}{
		"Authorization": getMemberToken,
	}

	//member api call to fetch member informations
	response, err := utils.APIRequest(http.MethodGet, cfg.MemberServiceURL+"/members/oauth", header, body)
	if err != nil {
		log.Errorf("OauthLogIn contoller-login get member api request failed %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"response": "token failure due to member service down"})
		return
	}
	if response.StatusCode == http.StatusBadRequest {
		log.Errorf("OauthLogIn contoller-member info fetch api error: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"response": "user not found,please register"})
		return
	}

	if response.StatusCode == http.StatusInternalServerError {
		log.Errorf("OauthLogIn contoller-member info fetch api error: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"response": "unable to fetch member informations"})
		return
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("OauthLogIn contoller-read member api response body error :%v", err)
	}

	if err := json.Unmarshal(responseBody, &apiResponseLogin); err != nil {
		log.Errorf("OauthLogIn contoller-failed to unmarshal member api response error :%v", err)
	}

	jwtPayload := entities.OAuthData{}
	jwtPayload.MemberEmail = apiResponseLogin.Data.Email
	memberID := apiResponseLogin.Data.MemberID.String()
	jwtPayload.MemberID = &memberID
	jwtPayload.PartnerID = apiResponseLogin.Data.PartnerID.String()
	jwtPayload.MemberName = apiResponseLogin.Data.Name
	jwtPayload.MemberType = apiResponseLogin.Data.MemberType
	jwtPayload.Roles = apiResponseLogin.Data.MemberRoles
	jwtPayload.PartnerName = apiResponseLogin.Data.Name

	tokenString := utilities.GenerateJwtToken(jwtPayload, consts.ExpTime, cfg.JwtKey)
	refreshToken := utilities.GenerateJwtToken(jwtPayload, consts.RefExpTime, cfg.JwtKey)
	refTok.RefreshToken = refreshToken

	err = oauth.useCase.PostRefreshToken(context.Background(), refTok, tokenString, partnerID, jwtPayload.MemberID)
	if err != nil {
		log.Errorf("OauthLogIn contoller-token entry failed in refresh table error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"response": "token failure due to member service down"})
		return
	}

	responseToken := entities.Response{
		Error:   nil,
		Message: "OAuth internal callback successful",
		Data: []map[string]interface{}{
			{
				"token":        tokenString,
				"refreshToken": refreshToken,
			},
		},
	}

	log.Printf("oauth internal callback successful")

	ctx.JSON(http.StatusOK, gin.H{
		"data": responseToken})

}
