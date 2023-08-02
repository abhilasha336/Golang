package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	facebookOAuth "golang.org/x/oauth2/facebook"
	googleOAuth "golang.org/x/oauth2/google"
	"golang.org/x/oauth2/spotify"
)

var strCsrf = "abhilash"
var jwtKey = []byte("jwtkey")

var (
	facebookClientID     = os.Getenv("FACEBOOK_CLIENT_ID")
	facebookClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	facebookRedirectURL  = os.Getenv("FACEBOOK_REDIRECT_URL")

	googleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL  = os.Getenv("GOOGLE_REDIRECT_URL")
)

var (
	spotifyClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	spotifyRedirectURL  = os.Getenv("SPOTIFY_REDIRECT_URL")
)

type MyData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Hd            string `json:"hd"`
}
type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type Claims struct {
	UserId               string // Change the data type of UserId to match your actual data type
	Email                string
	jwt.RegisteredClaims // Embedded to include the standard JWT claims
}

var spotifyOAuthConfig = &oauth2.Config{
	RedirectURL:  spotifyRedirectURL,
	Scopes:       []string{"user-read-email"},
	ClientID:     spotifyClientID,
	ClientSecret: spotifyClientSecret,
	Endpoint:     spotify.Endpoint,
}

var facebookOAuthConfig = &oauth2.Config{
	ClientID:     facebookClientID,
	ClientSecret: facebookClientSecret,
	RedirectURL:  facebookRedirectURL,
	Endpoint:     facebookOAuth.Endpoint,
	Scopes:       []string{"email", "public_profile"},
}

var googleOAuthConfig = &oauth2.Config{
	ClientID:     googleClientID,
	ClientSecret: googleClientSecret,
	RedirectURL:  googleRedirectURL,
	Endpoint:     googleOAuth.Endpoint,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
}

func handlerSpotifyLogin(c *gin.Context) {
	url := spotifyOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func handlerSpotifyCallback(c *gin.Context) {
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

	var user MyData
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error decoding user details: %v", err))
		return
	}
	tokenString := GenerateJwtToken(user, 1)
	fmt.Fprint(c.Writer, user)
	fmt.Fprint(c.Writer, tokenString)
}

func main() {
	r := gin.Default()

	r.GET("/", handlerHome)
	r.GET("/oauth", oauth)
	r.GET("/login/facebook", handlerFacebookLogin)
	r.GET("/login/google", handlerGoogleLogin)
	r.GET("/callback", handlerFacebookCallback)
	r.GET("/google/callback", handlerGoogleCallback)
	r.GET("/login/spotify", handlerSpotifyLogin)
	r.GET("/spotify/callback", handlerSpotifyCallback)

	if err := r.Run(":3000"); err != nil {
		log.Println("Error running server:", err)
	}
}

func oauth(c *gin.Context) {
	oauthType := c.DefaultQuery("type", "")
	switch oauthType {
	case "google":
		url := googleOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, url)
	case "facebook":
		url := facebookOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, url)
	case "spotify":
		url := spotifyOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, url)
	default:
		c.String(http.StatusBadRequest, "Invalid OAuth type")
	}
}

func handlerHome(c *gin.Context) {
	html := `<html>
		<body>
			<a href="/login/facebook">Facebook Sign in</a>
			<br>
			<a href="/login/google">Google Sign in</a>
			<br>
			<a href="/login/spotify">Spotify</a>

		</body>
	</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func handlerFacebookLogin(c *gin.Context) {
	url := facebookOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func handlerGoogleLogin(c *gin.Context) {
	url := googleOAuthConfig.AuthCodeURL(strCsrf, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func handlerFacebookCallback(c *gin.Context) {
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

	var data MyData
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	tokenString := GenerateJwtToken(data, 1)
	fmt.Fprint(c.Writer, tokenString)
}

func handlerGoogleCallback(c *gin.Context) {
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

	var data MyData
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 2nd functionalities
	tokenString := GenerateJwtToken(data, 1)

	fmt.Fprint(c.Writer, tokenString)
}

// generated a new JWT token from the userid with given expiry time
func GenerateJwtToken(data MyData, expTime int) string {
	expirationTime := time.Now().Add(time.Duration(expTime) * time.Minute)
	newClaims := &Claims{
		UserId: data.ID,
		Email:  data.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	token, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return token
}
