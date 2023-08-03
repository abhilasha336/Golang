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
	facebookOAuth "golang.org/x/oauth2/facebook"
)

var (
	facebookClientID     = os.Getenv("FACEBOOK_CLIENT_ID")
	facebookClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	facebookRedirectURL  = os.Getenv("FACEBOOK_REDIRECT_URL")
)

var facebookOAuthConfig = &oauth2.Config{
	ClientID:     facebookClientID,
	ClientSecret: facebookClientSecret,
	RedirectURL:  facebookRedirectURL,
	Endpoint:     facebookOAuth.Endpoint,
	Scopes:       []string{"email", "public_profile"},
}

type OauthFacebookController struct {
	router  *gin.RouterGroup
	useCase usecases.OuathFacebookUsecaseImply
}

func NewOauthFacebookController(router *gin.RouterGroup, useCase usecases.OuathFacebookUsecaseImply) *OauthFacebookController {
	return &OauthFacebookController{
		router:  router,
		useCase: useCase,
	}
}

func (oauth *OauthFacebookController) InitRoutes() {
	oauth.router.GET("/login/facebook", func(ctx *gin.Context) {
		// version.RenderHandler(ctx, oauth, "HandlerGoogleLogin")
		url := facebookOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
		ctx.Redirect(http.StatusFound, url)
	})

	oauth.router.GET("/callback", func(c *gin.Context) {
		// version.RenderHandler(c, oauth, "HandlerGoogleCallback")
		log.Println("hitted fb server")
		if c.Query("state") != strCsrf {
			fmt.Println("state is not valid")
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		token, err := facebookOAuthConfig.Exchange(context.Background(), c.Query("code"))
		if err != nil {
			fmt.Fprintln(c.Writer, err.Error())
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		facebookUserDetailsRequest, err := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email&access_token="+token.AccessToken, nil)
		if err != nil {
			fmt.Println("Error creating Facebook user details request:", err)
		}

		resp, err := http.DefaultClient.Do(facebookUserDetailsRequest)
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

		tokenString := utilities.GenerateJwtToken(data, 1)
		fmt.Fprint(c.Writer, tokenString)
	})
	// 	oauth.router.GET("/", func(ctx *gin.Context) {
	// 		// version.RenderHandler(ctx, oauth, "HandlerHome")
	// 		html := `<html>
	// 			<body>
	// 				<a href="/login/facebook">Facebook Sign in</a>
	// 				<br>
	// 				<a href="/login/google">Google Sign in</a>
	// 				<br>
	// 				<a href="/login/spotify">Spotify</a>

	//			</body>
	//		</html>`
	//		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	//	})
}
