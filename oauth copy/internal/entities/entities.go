package entities

import (
	"github.com/golang-jwt/jwt/v4"
)

// Response struct used to display error and data as response
type Response struct {
	Error   error
	Message string
	Data    []map[string]interface{}
}

// OAuthData used to store user data fetched from oauth servers
type OAuthData struct {
	ID            string      `json:"sub"`
	FID           string      `json:"id"`
	Email         string      `json:"email"`
	VerifiedEmail bool        `json:"email_verified"`
	Picture       string      `json:"picture"`
	GivenName     string      `json:"given_name"`
	FamilyName    string      `json:"family_name"`
	Locale        string      `json:"locale"`
	Name          string      `json:"display_name"`
	Country       string      `json:"country"`
	Type          string      `json:"type"`
	Product       string      `json:"product"`
	Followers     interface{} `json:"followers"`
	Hd            string      `json:"hd"`
	MemberID      *string     `json:"member_id"`
	PartnerID     string      `json:"partner_id"`
	PartnerName   string      `json:"partner_name"`
	MemberType    string      `json:"member_type"`
	Roles         []string    `json:"member_roles"`
	MemberEmail   string      `json:"member_email"`
	MemberName    string      `json:"member_name"`
}

// JwtValidateResponse used to store validated claims from jwttoken
type JwtValidateResponse struct {
	Valid       bool
	MemberID    *string  `json:"member_id"`
	MemberName  string   `json:"member_name"`
	PartnerID   string   `json:"partner_id"`
	PartnerName string   `json:"partner_name"`
	MemberType  string   `json:"member_type"`
	Roles       []string `json:"member_roles"`
	MemberEmail string   `json:"member_email"`
	ErrorMsg    string   `json:"errorms"`
}

// Claims used by jwt token to fetch claims
type Claims struct {
	MemberName  string
	MemberID    *string
	PartnerID   string
	PartnerName string
	MemberType  string
	Roles       []string
	MemberEmail string
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
