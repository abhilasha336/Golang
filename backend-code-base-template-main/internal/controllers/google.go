package controllers

import (
	"backend-code-base-template/internal/usecases"
	"backend-code-base-template/utilities"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	googleOAuth "golang.org/x/oauth2/google"
)

var strCsrf = "abhilash"
var (
	googleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL  = os.Getenv("GOOGLE_REDIRECT_URL")
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     googleClientID,
	ClientSecret: googleClientSecret,
	RedirectURL:  googleRedirectURL,
	Endpoint:     googleOAuth.Endpoint,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
}

type OauthGoogleController struct {
	router  *gin.RouterGroup
	useCase usecases.OuathGoogleUsecaseImply
}

func NewOauthGoogleController(router *gin.RouterGroup, useCase usecases.OuathGoogleUsecaseImply) *OauthGoogleController {
	return &OauthGoogleController{
		router:  router,
		useCase: useCase,
	}
}

func (oauth *OauthGoogleController) InitRoutes() {
	oauth.router.GET("/login/google", func(ctx *gin.Context) {
		// version.RenderHandler(ctx, oauth, "HandlerGoogleLogin")
		fmt.Println("hilogin")
		url := googleOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
		ctx.Redirect(http.StatusFound, url)
	})
	oauth.router.GET("/google/callback", func(c *gin.Context) {
		// version.RenderHandler(c, oauth, "HandlerGoogleCallback")
		log.Println("hitted google server")
		if c.Query("state") != strCsrf {
			fmt.Println("state is not valid")
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		token, err := googleOAuthConfig.Exchange(context.Background(), c.Query("code"))
		if err != nil {
			fmt.Fprintln(c.Writer, err.Error())
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		googleUserDetailsRequest, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo?access_token="+token.AccessToken, nil)
		if err != nil {
			fmt.Println("Error creating Google user details request:", err)
		}

		resp, err := http.DefaultClient.Do(googleUserDetailsRequest)
		if err != nil {
			fmt.Fprintln(c.Writer, err.Error())
			c.Redirect(http.StatusBadRequest, "/")
			return
		}
		defer resp.Body.Close()

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

		// 2nd functionalities
		tokenString := utilities.GenerateJwtToken(data, 1)

		fmt.Fprint(c.Writer, tokenString)
	})
	oauth.router.GET("/", func(ctx *gin.Context) {
		// version.RenderHandler(ctx, oauth, "HandlerHome")
		html := `<html>
			<body>
				<a href="/login/facebook">Facebook Sign in</a>
				<br>
				<a href="/login/google">Google Sign in</a>
				<br>
				<a href="/login/spotify">Spotify</a>

			</body>
		</html>`
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
}
