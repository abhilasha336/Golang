package controllers

import (
	"fmt"
	"net/http"
	"oauth/internal/consts"
	"oauth/internal/entities"
	"oauth/internal/usecases"
	"oauth/utilities"

	log "gitlab.com/tuneverse/toolkit/core/logger"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/core/version"
	"golang.org/x/oauth2"
)

// OauthController struct holds router group and usecase inetrface
type OauthController struct {
	router  *gin.RouterGroup
	useCase usecases.OuathUsecaseImply
	cfg     *entities.EnvConfig
}

// NewOauthController used to pass value of router and usecases
func NewOauthController(router *gin.RouterGroup, useCase usecases.OuathUsecaseImply, cfg *entities.EnvConfig) *OauthController {
	return &OauthController{
		router:  router,
		useCase: useCase,
		cfg:     cfg,
	}
}

// InitRoutes function used to init all routes
func (oauth *OauthController) InitRoutes() {

	oauth.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "HealthHandler")
	})

	oauth.router.GET("/:version/oauth/sso", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "OauthSso")
	})

	oauth.router.GET("/:version/partner", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "OauthHandler")
	})
	oauth.router.GET("/:version/callback", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "OauthCallBack")
	})
	oauth.router.PATCH("/:version/logout", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "OauthLogOut")
	})
	oauth.router.POST("/:version/refresh", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "OauthRefresh")
	})
	oauth.router.POST("/:version/login", func(ctx *gin.Context) {
		version.RenderHandler(ctx, oauth, "OauthLogIn")
	})

}

// OauthHandler function construct url hit to oauth google,fb,spotify respective servers
func (oauth *OauthController) OauthHandler(ctx *gin.Context) {

	var (
		err error
		log = log.Log().WithContext(ctx)
	)

	clientID := ctx.Request.Header.Get("client_id")
	client_Secret := ctx.Request.Header.Get("client_secret")
	provider := ctx.Request.Header.Get("provider")
	partnerID, redirectUri, err := oauth.useCase.GetPartnerId(ctx, clientID, client_Secret)
	if err != nil {
		log.Errorf("OauthHandler controller-partnerverify error:%v", err)
		return
	}

	if partnerID == "" || redirectUri == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"partner_id": "Invalid partner",
		})
		return
	}

	oauthData, err := oauth.useCase.GetOauthCredentials(ctx, provider, partnerID)
	if err != nil {
		log.Errorf("error in loading oauth credentials %v", err)
		response := entities.Response{
			Error:   fmt.Errorf("unable to load oauth credentials"),
			Message: "failed to load oauth credentials",
			Data:    nil,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	switch provider {
	case consts.GoogleProvider:
		entities.Endpoint = entities.Endpoints[consts.GoogleProvider]
	case consts.FacebookProvider:
		entities.Endpoint = entities.Endpoints[consts.FacebookProvider]
	case consts.SpotifyProvider:
		entities.Endpoint = entities.Endpoints[consts.SpotifyProvider]
	}

	config := utilities.Config(oauthData)
	url := config.AuthCodeURL(oauthData.State, oauth2.AccessTypeOffline)
	log.Printf("redirected successfully")
	ctx.JSON(http.StatusOK, gin.H{"url": url})

}
