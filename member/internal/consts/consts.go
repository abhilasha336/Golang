package consts

import "errors"

// DatabaseType represents the type of the database, set to "postgres."
const DatabaseType = "postgres"

// AppName represents the name of the application, set to "member."
const AppName = "member"

// AcceptedVersions represents an accepted API version, set to "v1."
const AcceptedVersions = "v1"

// ContextAcceptedVersions represents the key for the accepted API versions in a context.
const ContextAcceptedVersions = "Accept-Version"

// ContextSystemAcceptedVersions represents the key for system-accepted API versions in a context.
const ContextSystemAcceptedVersions = "System-Accept-Versions"

// ContextAcceptedVersionIndex represents the key for the accepted API version index in a context.
const ContextAcceptedVersionIndex = "Accepted-Version-index"

// ContextErrorResponses represents the key for error responses in a context.
const ContextErrorResponses = "context-error-response"

// ContextLocallizationLanguage represents the key for localization language in a context.
const ContextLocallizationLanguage = "lan"

// HeaderLocallizationLanguage represents the key for the Accept-Language header.
const HeaderLocallizationLanguage = "Accept-Language"

// CacheErrorData represents the cache name for error data.
const CacheErrorData = "CACHE_ERROR_DATA"

// ExpiryTime represents the cache expiry time, set to 180 (seconds).
const ExpiryTime = 180

// ContextEndPoints represent
const ContextEndPoints = "context-endpoints"

// KeyNames
const (
	//Parse Error indicates some error in the json data trying to parse
	ParseErr = "Parse_error"

	// ValidationErr represents a validation error.
	ValidationErr = "validation_error"

	// ForbiddenErr represents a forbidden error.
	ForbiddenErr = "forbidden"

	// UnauthorisedErr represents an unauthorized error.
	UnauthorisedErr = "unauthorized"

	// NotFound represents a not found error.
	NotFoundErr = "not_found"

	// InternalServerErr represents an internal server error.
	InternalServerErr = "internal_server_error"

	// Errors is a general category for error messages.
	Errors = "errors"

	// AllError represents all types of errors.
	AllError = "AllError"

	// Registration is related to user registration.
	Registration = "registration"

	// ErrorCode represents an error code.
	ErrorCode = "errorCode"

	// MemberIDErr represents an error related to "member_id".
	MemberIDErr = "member_id"

	// Required represents that a field is required.
	Required = "required"

	// FirstName represents a first name.
	FirstName = "firstname"

	// Invalid represents that a value is invalid.
	Invalid = "invalid"

	// LastName represents a last name.
	LastName = "lastname"

	// Email represents an email address.
	Email = "email"

	// EmailExists represents an error for an existing email address.
	EmailExists = "email_exists"

	// NewPassword represents a new password.
	NewPassword = "new_password"

	// CurrentPassword represents the current password.
	CurrentPassword = "current_password"

	// Format represents an error related to the format of a value.
	Format = "format"

	// InvalidPassword represents an error for an invalid password.
	InvalidPassword = "invalid_password"

	// Incorrect represents an error for something being incorrect.
	Incorrect = "incorrect"

	// Country represents a country.
	Country = "country"

	// State represents a state or region.
	State = "state"

	// Address represents an address.
	Address  = "address"
	Address1 = "address1"
	Address2 = "address2"
	// City represents a city.
	City = "city"

	// Zipcode represents a ZIP code.
	Zipcode = "zipcode"

	// PhoneNumber represents a phone number.
	PhoneNumber = "phone"

	// Valid represents that a value is valid.
	Valid = "valid"

	//MinLength represents the minimum required length of password
	MinLength = "min"

	PostMemberRegistration = "post_member_registration"
	Length                 = "length"
	Exists                 = "exists"
	Password               = "password"
	MinLengthPassword      = "min_length"
	FormatUpperCase        = "format_uppercase"
	FormatLowerCase        = "format_lowercase"
	FormatSpecialCharacter = "format_specialcharacter"
	FormatSpace            = "format_space"
	TermsAndConditions     = "terms_and_conditions"
	PayingTax              = "paying_tax"
	GetMemberProfile       = "get_member_profile"
	ProviderName           = "provider"
	MaximumAddressCount    = "maximumcount"
	Limit                  = "limit"
	CountryNotExist        = "countrynotexist"
	StateNotExist          = "statenotexist"
	BillingAddress         = "billingaddress"
	AlreadyExists          = "alreadyexists"
	Primary                = "primary"
	HasPrimary             = "hasprimary"
	ZipFormat              = "zipformat"
	MemberID               = "memberid"
	ChangeToFalse          = "changetofalse"
	MemberBilling          = "memberbilling"
	InvalidRelation        = "invalidrelation"
	PrimaryMandatory       = "primarymandatory"
	BillingCount           = "count"
	NoBillingAddress       = "noaddress"
	Title                  = "title"
	UpdateZipcode          = "updatezip"
	ExceedsLimit           = "Exceeds maximum record limit"
	UpdateStateWithCountry = "updatestate"
	InvalidKey             = "invalidkey"
	MaximumAddressLength   = "maximum"
	InvalidEmail           = "invalid_email"
	InValidProvider        = "invalid_provider"
	NoRelation             = "no_relation"
	TooLong                = "too_long"
	SuccessDelete          = "Successfully deleted Billing Address"
	NoMoreBillingAddress   = "no_address"
	MaxInt                 = 2147483647
	CannotCancel           = "cannot_cancel"
	AboutToExpire          = "about_to_expire"
	InGrace                = "in_grace"
	Expired                = "expired"
	PartnerID              = "partner_id"
	AuthenticationFailed   = "auth_failed"
	SuccessfullyDeleted    = "Successfully deleted member profile"
	Deleted                = "deleted"
)
const (
	// EndpointErr represents an error message related to loading endpoints from a service.
	EndpointErr = "Error occurred while loading endpoints from service"

	// ContextErr represents an error message related to loading an error from a service.
	ContextErr = "Error occurred while loading error from service"

	// Success is a success message for changing a password.
	Success = "Successfully changed password"

	// Successful is a success message for updating a profile.
	SuccessfullyUpdated = "Successfully Updated Profile"

	// SuccessfullyAdded is a success message for adding a billing address.
	SuccessfullyAdded = "Billing Address Added Successfully"

	// SuccessUpdated is a success message for updating a billing address.
	SuccessUpdated = "Billing Address Successfully Updated"

	//SuccessfullyListed is a success message for listing All Billing Address Succesfully
	SuccessfullyListed = "Billing Address Listed Successfully"

	// Mandatory represents a message for mandatory fields.
	Mandatory = "Mandatory Fields"

	//Minimum number of digits for zipcode
	Minimum = "minimum"

	// Active represents the status value for active members.
	Active = "active"
	// Inactive represents the status value for inactive members.
	Inactive = "inactive"
	// NotFound represents a not found error.
	NotFound = "not_found"
	// SuccessfullyRegistered is a message indicating successful member registration.
	SuccessfullyRegistered = "Member Registered successfully"

	// DefaultStatus is the default value for the "status" query parameter.
	DefaultStatus = "active"
	// DefaultPage is the default value for the "page" query parameter.
	DefaultPage = 1
	// DefaultLimit is the default value for the "limit" query parameter.
	DefaultLimit = 10
	//DefaultLimitAddress value for number of billing address listed per page .
	DefaultLimitAddress = 5
	//Maximum Allowed Limit value for limit parameter.
	MaximumLimit = 50
	// DefaultSortBy is the default value for the "sortby" query parameter.
	DefaultSortBy = "firstname"
	//DefaultOrder is the default order in which sort criteria works when nothing is specified explicitly.
	DefaultOrder = "ASC"
	// DefaultPartner is the default value for the "partner" query parameter.
	DefaultPartner = "-1"
	// DefaultRole is the default value for the "role" query parameter.
	DefaultRole = "-1"
	// DefaultSearch is the default value for the "search" query parameter.
	DefaultSearch = ""
	// MaximumNameLength is the maximum length for the name.
	MaximumNameLength = 63
	// ProviderInternal represents the internal authentication provider.
	ProviderInternal      = "internal"
	ContextPartnerID      = "partner_id"
	ProviderSpotify       = "spotify"
	Maximum               = "maximum"
	DefaultGender         = ""
	Gender                = "gender"
	ResetEmail            = "email"
	Key                   = "key"
	StoresExists          = "exists"
	DefaultCountry        = ""
	SuccessfullyInitiated = "Successfully Initiated Password Reset by Generating Key"
	Match                 = "failedmatch"
	// SubscriptionPlanRenewableDuration represents whether the subscription plan is renewable within certain duration.
	SubscriptionPlanRenewableDuration = "can_renewable_within"
	LimitReached                      = "limit"

	DefaultSortByID              = "id"
	MaximumRequestError          = "Cannot exceed maximum page limit"
	SuccessfullyAddedMemberStore = "Successfully Added New Stores"
	InvalidAddress               = "invalid_address"

	CheckStatus      = "check_status"
	NotSubscribedYet = "not_subscribed"
	NotEnabled       = "notenabled"

	Subscribed     = "alreadysubscribed"
	SubscribedOnce = "oncesubscribed"
	NoField        = "No fields to update"
	CustonName     = "name"
	Name           = "name"
	Empty          = "empty"
	Language       = "lang"
	InvalidFormat  = "invalid_format"
	NoPayin        = "no_payin"
	NoDetails      = "no_details"
)

// SuccessfullyCheckedout is a constant representing a success message for checking out a subscription plan.
const SuccessfullyCheckedout = "Successfully checked out subscription plan"

// SuccessfullyCancelled is a constant representing a success message for cancellation of  a subscription plan.
const SuccessfullyCancelled = "Successfully cancelled  subscription plan"

// SuccessfullyCheckedout is a constant representing a success message for checking out a subscription plan.
const SuccessfullyRenewed = "Successfully Renewed existing Subscription"

// SubscriptionID of the plan
const SubscriptionID = "subscription_id"

// PaymentGatewayID of the member
const PaymentGatewayID = "payment_gateway_id"

// DummyPaymentGatewayURL  a dummy url to redirect to payment gateway
const DummyPaymentGatewayURL = "https://dummy-payment-gateway.com/pay"

var ErrInParsing = errors.New("Error occured while parsing")

const (
	// Page represents the page parameter for pagination.
	Page = "page"
	// Limit represents the limit parameter for pagination.
)

// KeyNames
const (
	// SubscriptionPlan represents a subscription plan key.
	SubscriptionPlan = "subscription_plan"
)

// KeyNames
const (
	SuccessSwitched         = "Successfully switched products between subscriptions"
	SuccessfullyListedPlans = "Subscription plans listed successfully"
	LimitExceeds            = "limit_exceeds"
	NotActive               = "not_active"
	NotAllowed              = "not_allowed"
	AlreadyExist            = "already_exists"
)

// success response code
const (
	StatusOk = 200
)

// validation errors
const (
	MemberrID             = "member_id"
	CurrentSubscriptionID = "current_subscription_id"
	NewSubscriptionID     = "new_subscription_id"
	ProductReferenceID    = "product_reference_id"
	CurrentStatus         = "current_status"
	PaymentFailed         = "payment_failed"
	Processing            = "processing"
	OnHold                = "on_hold"
	Cancelled             = "cancelled"
)

// consts used in utilities
const (
	LimitVal  = 25
	InputOne  = 10
	InputTwo  = 64
	RateLimit = 50
)

// Cache Keys
const (
	CacheErrorKey     = "ERROR_CACHE_KEY_LABEL"
	CacheEndpointsKey = "endpoints"
)

// consts used in app
const (
	LogMaxAge    = 7
	LogMaxSize   = (1024 * 1024 * 10)
	LogMaxBackup = 5
)

const (
	// DefaultLimit specifies the default limit for pagination.
	LimitDefault = 10
	// MaxAllowedLimit specifies the maximum allowed limit for pagination.
	MaxAllowedLimit = 25
)
