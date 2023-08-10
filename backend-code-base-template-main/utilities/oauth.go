package utilities

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("jwtkey")

type MyData struct {
	ID            string `json:"sub"`
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	Hd            string `json:"hd"`
}
type JwtValidateResponse struct {
	Valid    bool
	UserId   string `json:"userid"`
	EmailId  string `json:"emailid"`
	ErrorMsg string `json:"errorms"`
}

type Claims struct {
	UserId               string // Change the data type of UserId to match your actual data type
	Email                string
	jwt.RegisteredClaims // Embedded to include the standard JWT claims
}

func GenerateJwtToken(data MyData, expTime int) string {
	var newClaims *Claims
	expirationTime := time.Now().Add(time.Duration(expTime) * time.Minute)
	if data.ID == "" {
		newClaims = &Claims{
			UserId: data.Id,
			Email:  data.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
	} else {
		newClaims = &Claims{
			UserId: data.ID,
			Email:  data.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	token, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return token
}

func ValidateJwtToken(token string) (response JwtValidateResponse) {

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Print("panicked")
			response.ErrorMsg = "Token seems to be invalid!"
		}
	}()

	response = JwtValidateResponse{
		Valid:    false,
		UserId:   "0",
		EmailId:  "",
		ErrorMsg: "",
	}

	claims := &Claims{}
	if len(token) == 0 {
		response.ErrorMsg = "No authorization token passed"
	} else {
		tokenParsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if claims, ok := tokenParsed.Claims.(*Claims); ok && tokenParsed.Valid {
			response.Valid = true
			response.UserId = claims.UserId
			response.EmailId = claims.Email
		} else {
			response.ErrorMsg = err.Error()
		}
	}
	return response
}
