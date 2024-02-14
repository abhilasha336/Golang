package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/models"
)

// Member struct to hold details of member profile
type Member struct {
	Title                 string `json:"title"`
	FirstName             string `json:"firstname"`
	LastName              string `json:"lastname"`
	Email                 string `json:"email"`
	Gender                string `json:"gender"`
	Language              string `json:"language"`
	Country               string `json:"country"`
	State                 string `json:"state"`
	Address1              string `json:"address1"`
	Address2              string `json:"address2"`
	City                  string `json:"city"`
	Zipcode               string `json:"zip"`
	Phone                 string `json:"phone"`
	TermsConditionChecked bool   `json:"terms_condition_checked"`
	PayingTax             bool   `json:"paying_tax"`
	Password              string `json:"password,omitempty"`
	Provider              string `json:"provider,omitempty"`
}

// BillingAddress struct to hold details
type BillingAddress struct {
	Address string `json:"address" `
	Zipcode string `json:"zipcode" `
	Country string `json:"country"  `
	State   string `json:"state" `
	Primary bool   `json:"primary"`
}

// BillingAddressResponse response struct to hold details
type BillingAddressResponse struct {
	BillingAddress []BillingAddress `json:"billing_address"`
}

// PasswordChangeRequest represents a request to change a password.
type PasswordChangeRequest struct {
	Key             string `json:"key"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ErrorResponse represents an error response that can be sent back to clients in JSON format.
type ErrorResponse struct {
	Message   string                 `json:"message"`
	ErrorCode interface{}            `json:"errorCode"`
	Errors    map[string]interface{} `json:"errors"`
}

// MemberProfile combines MemberDetails and an array of BillingAddress
type MemberProfile struct {
	MemberDetails        Member           `json:"member_details"`
	MemberBillingAddress []BillingAddress `json:"member_billing_address"`
	EmailSubscribed      bool             `json:"email_subscribed"`
}

// Params represents a set of parameters that can be used for querying members.
type Params struct {
	Status  string
	Page    int16
	Limit   int16
	SortBy  string
	Partner string
	Role    string
	Search  string
	Order   string
	Country string
	Gender  string
}
type ResetPassword struct {
	Email string `json:"email"`
}

// Role represents information about a specific role with an ID and name.
type Role struct {
	Id   int
	Name string
}

// Country represents information about a specific country with its code and name.
type Country struct {
	Code string
	Name string
}

// ViewMembers represents information about members with nested Role and Country data.
type ViewMembers struct {
	MemberId    uuid.UUID
	Name        string
	Role        Role
	Gender      string
	PartnerName string
	Email       string
	Country     Country
	Active      bool
	AlbumCount  int
	TrackCount  int
	ArtistCount int
}

// BasicMemberData represents basic member details.
type BasicMemberData struct {
	MemberID    uuid.UUID `json:"member_id"`
	Name        string    `json:"member_name"`
	PartnerName string    `json:"partner_name"`
	PartnerID   uuid.UUID `json:"partner_id"`
	ProviderID  uuid.UUID `json:"provider_id"`
	Email       string    `json:"member_email"`
	MemberType  string    `json:"member_type"`
	MemberRoles []string  `json:"member_roles"`
}

// MemberPayload represents essential information about a member.
type MemberPayload struct {
	Email    string `json:"email"`
	Provider string `json:"provider"`
	Password string `json:"password"`
}

type JwtValidateResponse struct {
	Valid       bool
	MemberID    string   `json:"member_id"`
	MemberName  string   `json:"member_name"`
	PartnerID   string   `json:"partner_id"`
	PartnerName string   `json:"partner_name"`
	MemberType  string   `json:"member_type"`
	Roles       []string `json:"member_roles"`
	MemberEmail string   `json:"member_email"`
	ErrorMsg    string   `json:"errorms"`
}

type Claims struct {
	MemberName  string
	MemberID    string
	PartnerID   string
	PartnerName string
	MemberType  string
	Roles       []string
	MemberEmail string
	jwt.RegisteredClaims
}

// GetMemberByID struct to store details of member for operations like updation
type MemberByID struct {
	Title     sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString
	Country   sql.NullString
	State     sql.NullString
	Zipcode   sql.NullString
	Phone     sql.NullString
	City      sql.NullString
	Address1  sql.NullString
	Address2  sql.NullString
}

// CheckoutSubscription represents the data structure for a subscription checkout request.
type CheckoutSubscription struct {
	SubscriptionID   string `json:"subscription_id"`    // SubscriptionID is the unique identifier for the subscription.
	PaymentGatewayID int    `json:"payment_gateway_id"` // PaymentGatewayID is the ID of the payment gateway to use for the subscription payment.
	CustomName       string `json:"custom_name"`        //CustomName is the alternative name for the subscribed plan
}

// SubscriptionRenewal represents the data structure for a subscription renewal request.
type SubscriptionRenewal struct {
	MemberSubscriptionID string `json:"member_subscription_id"` // SubscriptionID is the unique identifier for the subscription.
	PaymentGatewayID     int    `json:"payment_gateway_id"`     // PaymentGatewayID is the ID of the payment gateway to use for the subscription payment.

}

// CancelSubscription structure
type CancelSubscription struct {
	MemberSubscriptionID string `json:"member_subscription_id"` // SubscriptionID is the unique identifier for the subscription.
}

// SwitchSubscriptions represents the SwitchSubscriptions entity.
type SwitchSubscriptions struct {
	CurrentSubscriptionID string `json:"current_subscription_id"`
	NewSubscriptionID     string `json:"new_subscription_id"`
	ProductReferenceID    string `json:"product_reference_id"`
}

// ListAllSubscriptions represents the ListAllSubscriptions entity.
type ListAllSubscriptions struct {
	ID                  uuid.UUID           `json:"id,omitempty"`
	CustomNameJSON      json.RawMessage     `json:"custom_name,omitempty"`
	CustomName          NullableString      `json:"-"`
	Status              string              `json:"status,omitempty"`
	ExpirationDate      string              `json:"expiration_date"`
	ProductsAdded       int                 `json:"products_added"`
	TracksAdded         int                 `json:"tracks_added"`
	ArtistsAdded        int                 `json:"artists_added"`
	SubscriptionDetails SubscriptionDetails `json:"subscription_details"`
}

// SubscriptionDetails represents subscription entity
type SubscriptionDetails struct {
	SubscriptionID  string `json:"subscription_id,omitempty"`
	Name            string `json:"name,omitempty"`
	SKU             string `json:"sku,omitempty"`
	Duration        string `json:"duration,omitempty"`
	MaximumProducts int    `json:"maximum_products"`
	MaximumTracks   int    `json:"maximum_tracks"`
	MaximumArtists  int    `json:"maximum_artists"`
}

// MetaData represents metadata for an entity.
type MetaData struct {
	Total       int64 `json:"total"`
	PerPage     int32 `json:"per_page"`
	CurrentPage int32 `json:"current_page"`
	Next        int32 `json:"next"`
	Prev        int32 `json:"prev"`
}

// SuccessResponse represents success data for an entity.
type SuccessResponse struct {
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Metadata interface{}            `json:"metadata"`
	Data     []ListAllSubscriptions `json:"data"`
}

// ReqParams represents the request parameters.
type ReqParams struct {
	Page   int32  `form:"page"`
	Limit  int32  `form:"limit"`
	Search string `form:"search"`
	Sort   string `form:"sort"`
	Status string `form:"status"`
}

type NullableString struct {
	sql.NullString
}

// SubscriptionDetails represent curent deails about subscription
type SubscriptionStatusDetails struct {
	IsInGracePeriod       bool
	GraceStart            time.Time
	GraceEnd              time.Time
	GraceDuration         time.Duration
	IsWithinGraceDuration bool
	Err                   error
}
type AddMemberStores struct {
	Storelist []string `json:"stores"`
}
type Response struct {
	BillingAddress []BillingAddress `json:"billing_address"`
	Message        string           `json:"message"`
}

// PaymentGatewayDetails,Define a struct to represent the JSON data
type PaymentGatewayDetails struct {
	Gateway               string   `json:"gateway"`
	Email                 string   `json:"email"`
	ClientID              string   `json:"client_id"`
	ClientSecret          string   `json:"client_secret"`
	Payin                 bool     `json:"payin"`
	Payout                bool     `json:"payout"`
	SupportedCurrency     []string `json:"supported_currency"`
	DefaultPayinCurrency  string   `json:"default_payin_currency"`
	DefaultPayoutCurrency string   `json:"default_payout_currency"`
}
type MemberResponse struct {
	Status   string       `json:"status"`
	Code     int          `json:"code"`
	Message  string       `json:"message"`
	DataResp DataResponse `json:"data"`
}

type DataResponse struct {
	Metadata models.MetaData `json:"meta_data"`
	Data     []ViewMembers   `json:"member_data"`
}
type FailureResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}
