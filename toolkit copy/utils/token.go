package utils

import (
	"reflect"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTAuthToken(jwtTokenKey string, jwtTokenPayload interface{}) (string, error) {
	var claimsJWT jwt.Claims

	if reflect.TypeOf(jwtTokenPayload).Kind() == reflect.Map {
		// Create the Claims
		atClaims := jwt.MapClaims{}
		for k, v := range jwtTokenPayload.(map[string]interface{}) {
			atClaims[k] = v
		}
		claimsJWT = atClaims
	} else {
		claimsJWT = jwt.StandardClaims{}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsJWT)
	tokenString, err := token.SignedString([]byte(jwtTokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
