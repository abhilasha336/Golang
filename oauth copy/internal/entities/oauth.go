package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/spotify"
)

// Refresh which holds refreshtoken

type FormEmail struct {
	Email string `json:"email"`
}
type Refresh struct {
	RefreshToken string `json:"refreshtoken"`
}

// StringArray holds slice of string

// OAuthCredentials used to hold credentials fetched from database
type OAuthCredentials struct {
	ProviderName string      `json:"provider_name" db:"provider_name"`
	ClientID     string      `json:"client_id" db:"client_id"`
	ClientSecret string      `json:"client_secret" db:"client_secret"`
	RedirectURL  string      `json:"redirect_url" db:"redirect_url"`
	Scopes       StringArray `json:"scopes" db:"scopes"`
	State        string      `json:"state" db:"state"`
	TokenURL     string      `json:"token_url" db:"token_url"`
	JWTKey       string      `json:"jwt_key" db:"jwt_key"`
}

// Endpoint holds oauth's endpoints
var Endpoint oauth2.Endpoint

// Endpoints map variable to fetch endpoints of respective different oauth
var Endpoints = map[string]oauth2.Endpoint{
	"google":   google.Endpoint,
	"facebook": facebook.Endpoint,
	"spotify":  spotify.Endpoint,
}

type StringArray []string

// Scan unmarshal any values
func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}

	return json.Unmarshal(b, a)
}

// Value unmarshal any values
func (a *StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}
