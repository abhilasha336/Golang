package controllers

import (
	"backend-code-base-template/internal/usecases"
	"backend-code-base-template/utilities"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var (
	spotifyClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURL  = os.Getenv("SPOTIFY_REDIRECT_URL")
)

var spotifyOAuthConfig = &oauth2.Config{
	RedirectURL:  spotifyRedirectURL,
	Scopes:       []string{"user-read-email"},
	ClientID:     spotifyClientID,
	ClientSecret: spotifyClientSecret,
	Endpoint:     spotify.Endpoint,
}

type OauthSpotifyController struct {
	router  *gin.RouterGroup
	useCase usecases.OuathSpotifyUsecaseImply
}

func NewOauthSpotifyController(router *gin.RouterGroup, useCase usecases.OuathSpotifyUsecaseImply) *OauthSpotifyController {
	return &OauthSpotifyController{
		router:  router,
		useCase: useCase,
	}
}

func (oauth *OauthSpotifyController) InitRoutes() {
	oauth.router.GET("/login/spotify", func(ctx *gin.Context) {
		// version.RenderHandler(ctx, oauth, "HandlerGoogleLogin")
		fmt.Println("hilogin")
		url := spotifyOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
		ctx.Redirect(http.StatusFound, url)
	})
	oauth.router.GET("/spotify/callback", func(c *gin.Context) {
		log.Println("hitted Spotify server")
		if c.Query("state") != strCsrf {
			fmt.Println("state is not valid")
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		token, err := spotifyOAuthConfig.Exchange(c, c.Query("code"))
		if err != nil {
			fmt.Fprintln(c.Writer, err.Error())
			c.Redirect(http.StatusBadRequest, "/")
			return
		}
		client := spotifyOAuthConfig.Client(c, token)
		resp, err := client.Get("https://api.spotify.com/v1/me")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching user details: %v", err))
			return
		}

		// var user utilities.MyData
		// if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		// 	c.String(http.StatusInternalServerError, fmt.Sprintf("Error decoding user details: %v", err))
		// 	return
		// }
		// fmt.Printf("%+v", user)

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(c.Writer, err.Error())
			c.Redirect(http.StatusBadRequest, "/")
			return
		}
		fmt.Fprint(c.Writer, string(content))

		var data utilities.MyData
		err = json.Unmarshal([]byte(content), &data)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
		fmt.Printf("%+v", data)

		tokenString := utilities.GenerateJwtToken(data, 1)
		fmt.Fprint(c.Writer, data)
		fmt.Fprint(c.Writer, tokenString)

		res := utilities.ValidateJwtToken(tokenString)
		fmt.Fprintf(c.Writer, "%+v", res)

	})

}
