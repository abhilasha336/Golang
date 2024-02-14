package utilities

import (
	"oauth/internal/entities"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

// GenerateJwtToken function used to generate jwt token
func GenerateJwtToken(data entities.OAuthData, expTime int, jwtKey string) string {

	var newClaims *entities.Claims
	jwtKeyVal := []byte(jwtKey)
	expirationTime := time.Now().Add(time.Duration(expTime) * time.Minute)

	newClaims = &entities.Claims{
		MemberName:  data.MemberName,
		MemberID:    data.MemberID,
		PartnerID:   data.PartnerID,
		PartnerName: data.PartnerName,
		MemberType:  data.MemberType,
		Roles:       data.Roles,
		MemberEmail: data.MemberEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	token, err := jwtToken.SignedString(jwtKeyVal)
	if err != nil {
		return err.Error()
	}
	return token
}

// ValidateJwtToken function used to validate jwt token
func ValidateJwtToken(token string, jwtKey string) (response entities.JwtValidateResponse) {

	jwtKeyVal := []byte(jwtKey)

	defer func() {
		if rec := recover(); rec != nil {
			response.ErrorMsg = "Token seems to be invalid!"
		}
	}()

	response = entities.JwtValidateResponse{
		Valid:       false,
		MemberID:    nil,
		PartnerID:   "",
		PartnerName: "",
		MemberName:  "",
		Roles:       []string{},
		MemberEmail: "",
		ErrorMsg:    "",
	}

	claims := &entities.Claims{}
	if len(token) == 0 {
		response.ErrorMsg = "No authorization token passed"
	} else {
		tokenParsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKeyVal, nil
		})
		if claims, ok := tokenParsed.Claims.(*entities.Claims); ok && tokenParsed.Valid {
			response.Valid = true
			response.MemberID = claims.MemberID
			response.PartnerID = claims.PartnerID
			response.MemberType = claims.MemberType
			response.Roles = claims.Roles
			response.MemberEmail = claims.MemberEmail
			response.MemberName = claims.MemberName
			response.PartnerName = claims.PartnerName
		} else {
			response.ErrorMsg = err.Error()
		}
	}
	return response
}

// Config function used to load credential
func Config(credData entities.OAuthCredentials) *oauth2.Config {

	return &oauth2.Config{
		ClientID:     credData.ClientID,
		ClientSecret: credData.ClientSecret,
		RedirectURL:  credData.RedirectURL,
		Endpoint:     entities.Endpoint,
		Scopes:       credData.Scopes,
	}
}
