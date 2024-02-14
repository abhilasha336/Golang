package controllers

import (
	"context"
	"net/http"
	"oauth/internal/consts"
	"oauth/utilities"

	"github.com/gin-gonic/gin"
	log "gitlab.com/tuneverse/toolkit/core/logger"
)

// OauthCallBack function which fetch data from server with auth code and generate token
// dummy function, in realtime this functionality is provided by frontend
func (oauth *OauthController) OauthCallBack(ctx *gin.Context) {

	var (
		partnerID, provider string
		log                 = log.Log().WithContext(ctx)
	)

	provider = "google"
	partnerID = "614608f2-6538-4733-aded-96f902007254"

	oauthData, err := oauth.useCase.GetOauthCredentials(ctx, provider, partnerID)
	if err != nil {
		log.Errorf("OauthCallBack controller-error in loading oauth credentials error:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"response": "callback token genration failed"})
		return
	}

	config := utilities.Config(oauthData)

	if ctx.Query(consts.State) != oauthData.State {
		log.Errorf("OauthCallBack controller-state is not valid")
		ctx.Redirect(http.StatusBadRequest, "/")
		return
	}

	token, err := config.Exchange(context.Background(), ctx.Query("code"))
	if err != nil {
		log.Errorf("OauthCallBack controller-error in generating token %v", err)
		ctx.Redirect(http.StatusBadRequest, "/")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"react-token": token})

}
