package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"oauth/internal/consts"
	"oauth/internal/entities"
	"oauth/utilities"

	"github.com/gin-gonic/gin"
	log "gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/utils"
)

// fn handles sso ,communication with member post and get services and generates token
func (oauth *OauthController) OauthSso(ctx *gin.Context) {
	var (
		data                entities.OAuthData
		partnerID, provider string
		apiResponse         entities.BasicMemberDataResponse
		refTok, refToken    entities.Refresh
		log                 = log.Log().WithContext(ctx)
		cfg                 = oauth.cfg
		header, body        map[string]interface{}
	)

	provider = ctx.GetHeader("provider")
	partnerID = ctx.GetHeader("partner_id")

	oauthData, err := oauth.useCase.GetOauthCredentials(ctx, provider, partnerID)
	if err != nil {
		log.Errorf("OauthSso controller-error in loading oauth credentials %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"response": "unauthorised credentials"})
		return
	}
	config := utilities.Config(oauthData)

	token, err := utilities.GetTokenFromHeader(ctx.Request)
	if err != nil {
		log.Errorf("OauthSso controller-token conversion invalid:%v", err)
	}

	switch provider {
	case consts.SpotifyProvider:
		client := config.Client(ctx, token)
		resp, err := client.Get(oauthData.TokenURL)
		if err != nil {
			log.Errorf("OauthSso controller-error in getting spotify credentials: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"response": "token request spotify external server down"})
			return
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("OauthSso controller-error in reading oauth servers response body %v", err)
			ctx.Redirect(http.StatusBadRequest, "/")
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Errorf("OauthSso controller-unmarshal error: %v", err)
		}

	default:
		userDetailsRequest, err := http.NewRequest(http.MethodGet, oauthData.TokenURL+token.AccessToken, nil)
		if err != nil {
			log.Errorf("OauthSso controller-error in creating user details request %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"response": "token request external server down"})
			return
		}

		resp, err := http.DefaultClient.Do(userDetailsRequest)
		if err != nil {
			log.Errorf("OauthSso controller-error in fetching %v", err)
			ctx.Redirect(http.StatusBadRequest, "/")
			return
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("OauthSso controller-error in reading oauth servers response body %v", err)
			ctx.Redirect(http.StatusBadRequest, "/")
			return
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Errorf("OauthSso controller-unmarshal error: %v", err)
		}

	}

	if provider == "spotify" {
		body = map[string]interface{}{
			"email":    data.FID,
			"provider": provider,
		}
	} else {
		body = map[string]interface{}{
			"email":    data.Email,
			"provider": provider,
		}
	}
	header = map[string]interface{}{
		"partner_id": partnerID,
	}

	//api call to register member
	emailPostResponse, err := utils.APIRequest(http.MethodPost, cfg.MemberServiceURL+"/members", header, body)
	if emailPostResponse.StatusCode == http.StatusInternalServerError {
		log.Errorf("OauthSso controller-member register api failed: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"response": "member registration failed"})
		return
	}

	if emailPostResponse.StatusCode == http.StatusOK {
		log.Printf("OauthSso controller-member added successfully(new user via oAuth)")
	}
	if emailPostResponse.StatusCode == http.StatusBadRequest {
		log.Printf("OauthSso controller-email exist as a member")
	}

	//payload for fetching member info
	if provider == "spotify" {
		body = map[string]interface{}{
			"email":    data.FID,
			"provider": provider,
		}
	} else {
		body = map[string]interface{}{
			"email":    data.Email,
			"provider": provider,
		}
	}

	data.PartnerID = partnerID
	memberGetToken := utilities.GenerateJwtToken(data, consts.TempExpTime, cfg.JwtKey)
	refToken.RefreshToken = memberGetToken
	err = oauth.useCase.PostRefreshToken(context.Background(), refToken, memberGetToken, partnerID, data.MemberID)
	if err != nil {
		log.Errorf("OauthSso controller-dummy token entry in refreshtoken table failed: %v", err)
	}
	header = map[string]interface{}{
		"Authorization": memberGetToken,
	}

	memberInfoResponse, err := utils.APIRequest(http.MethodGet, cfg.MemberServiceURL+"/members/oauth", header, body)
	if err != nil {
		log.Error("OauthSso controller-get member service call failed:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"response": "token failure due to member service down"})
		return
	}
	responsebody, err := io.ReadAll(memberInfoResponse.Body)
	if err != nil {
		log.Errorf("unable to read member response body:%v", err)
	}

	// Unmarshal the response body into the struct
	if err := json.Unmarshal(responsebody, &apiResponse); err != nil {
		log.Errorf("unable to unmarshal member response body:%v", err)
	}

	memberProviderName, err := oauth.useCase.GetProviderName(ctx, (apiResponse.Data.ProviderID).String())
	if err != nil {
		log.Errorf("OauthSso controller-unable to fectch provider name using member api%v", err)
	}

	if memberProviderName != provider {
		log.Errorf("OauthSso controller-member already exist with another oauthprovider:%v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"response": "this email already registered with another oauthprovider:" + memberProviderName})
		return
	}

	jwtPayload := entities.OAuthData{}
	memberID := apiResponse.Data.MemberID.String()
	jwtPayload.MemberID = &memberID
	jwtPayload.PartnerID = apiResponse.Data.PartnerID.String()
	jwtPayload.MemberEmail = apiResponse.Data.Email
	jwtPayload.MemberName = apiResponse.Data.Name
	jwtPayload.MemberType = apiResponse.Data.MemberType
	jwtPayload.Roles = apiResponse.Data.MemberRoles
	jwtPayload.PartnerName = apiResponse.Data.Name
	tokenString := utilities.GenerateJwtToken(jwtPayload, consts.ExpTime, cfg.JwtKey)
	refreshToken := utilities.GenerateJwtToken(jwtPayload, consts.RefExpTime, cfg.JwtKey)
	refTok.RefreshToken = refreshToken

	err = oauth.useCase.PostRefreshToken(context.Background(), refTok, tokenString, partnerID, jwtPayload.MemberID)
	if err != nil {
		log.Printf("OauthSso controller-token entry in refreshtoken table failed: %v", err)
		responseFail := entities.Response{
			Error:   err,
			Message: "db refresh token failure",
			Data: []map[string]interface{}{
				{
					"message": "token failure",
				},
			},
		}
		ctx.JSON(http.StatusBadRequest, responseFail)
		return
	}

	response := entities.Response{
		Error:   nil,
		Message: "OAuth callback successful",
		Data: []map[string]interface{}{
			{
				"token":        tokenString,
				"refreshToken": refreshToken,
			},
		},
	}

	log.Printf("oauth callback successful")
	ctx.JSON(http.StatusOK, response)
}
