package utilities

import (
	"net/http"

	"golang.org/x/oauth2"
)

// fn converts token from request header to oauthtoken type
func GetTokenFromHeader(r *http.Request) (*oauth2.Token, error) {
	tokenString := r.Header.Get("Authorization") // Extract token string from the header

	// Create an OAuth2 Token instance with the given access token
	token := &oauth2.Token{
		AccessToken: tokenString,
		TokenType:   "Bearer", // Token type (e.g., Bearer)
	}

	return token, nil
}
