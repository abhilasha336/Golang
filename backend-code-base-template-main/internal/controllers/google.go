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

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// var strCsrf = "abhilash"

// var (
// 	googleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
// 	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
// 	googleRedirectURL  = os.Getenv("GOOGLE_REDIRECT_URL")
// )

type OAuthData struct {
	ProviderName string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	State        string
	TokenURL     string
	JWTKey       string
}

// var googleOAuthConfig = &oauth2.Config{
// 	ClientID:     googleClientID,
// 	ClientSecret: googleClientSecret,
// 	RedirectURL:  googleRedirectURL,
// 	Endpoint:     googleOAuth.Endpoint,
// 	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// }

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

var oauthData = OAuthData{
	ProviderName: "google",
	ClientID:     "859339718701-pvvlufnvtjq6b6rorc9h8q1ll3cjo2gp.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:3000/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	State:        "abhilash",
	TokenURL:     "https://www.googleapis.com/oauth2/v3/userinfo?access_token=",
	JWTKey:       "sampleJwtKey",
}

func Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oauthData.ClientID,
		ClientSecret: oauthData.ClientSecret,
		RedirectURL:  oauthData.RedirectURL,
		Endpoint:     endpoint,
		Scopes:       oauthData.Scopes,
	}
}

var endpoint oauth2.Endpoint
var endpoints = map[string]oauth2.Endpoint{
	"google":   google.Endpoint,
	"facebook": facebook.Endpoint,
	// Add more provider endpoints here
}

func (oauth *OauthGoogleController) InitRoutes() {
	oauth.router.GET("/login/google", func(ctx *gin.Context) {

		// version.RenderHandler(ctx, oauth, "HandlerGoogleLogin")
		g := "google"

		// endpoints = map[string]oauth2.Endpoint{
		// 	"google":   google.Endpoint,
		// 	"facebook": facebook.Endpoint,
		// 	// Add more provider endpoints here
		// }
		switch g {
		case "google":
			endpoint = endpoints["google"]
		}
		config := Config()
		fmt.Println("hilogin")
		url := config.AuthCodeURL(oauthData.State, oauth2.AccessTypeOffline)
		fmt.Println("urlll", url)
		ctx.Redirect(http.StatusFound, url)
	})
	oauth.router.GET("/google/callback", func(c *gin.Context) {
		// version.RenderHandler(c, oauth, "HandlerGoogleCallback")
		config := Config()

		log.Println("hitted google server")
		if c.Query("state") != oauthData.State {
			fmt.Println("state is not valid")
			c.Redirect(http.StatusBadRequest, "/")
			return
		}
		fmt.Println("stateeeeeeeee")
		token, err := config.Exchange(context.Background(), c.Query("code"))
		if err != nil {
			fmt.Fprintln(c.Writer, err.Error())
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		googleUserDetailsRequest, err := http.NewRequest("GET", oauthData.TokenURL+token.AccessToken, nil)
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
		fmt.Printf("%+v", data)

		// 2nd functionalities
		tokenString := utilities.GenerateJwtToken(data, 1)

		fmt.Fprint(c.Writer, tokenString)
		res := utilities.ValidateJwtToken(tokenString)
		fmt.Fprintf(c.Writer, "%+v", res)

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
