package utilities

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("jwtkey")

type MyData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Hd            string `json:"hd"`
}
type Claims struct {
	UserId               string // Change the data type of UserId to match your actual data type
	Email                string
	jwt.RegisteredClaims // Embedded to include the standard JWT claims
}

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
