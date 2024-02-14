package consts

const (

	// AppName is the name of the application.
	AppName = "localization"
	// AcceptedVersions contains the accepted versions.
	AcceptedVersions = "v1"
	DatabaseType     = "postgres"
)

// Context keys for accepted versions and language.
const (
	ContextAcceptedVersions       = "Accept-Version"
	ContextSystemAcceptedVersions = "System-Accept-Versions"
	ContextAcceptedVersionIndex   = "Accepted-Version-index"
	ContextLocallizationLanguage  = "Accept-Language"
)

// Error codes for localization.
const (
	CodeL001 = "L_001"
	CodeL200 = "L_200"
)

const (
	// CollectionStatusCodes contains status codes for collections.
	CollectionStatusCodes = "tuneverseNew"
	// CollectionEndpoint is the endpoint for collections.
	CollectionEndpoint = "endpointnames"
	// ListAllErr contains error messages for listing.
	ListAllErr = "AllError"
	// ValidataionErr contains error message for validation.
	ValidataionErr = "Validation Error"
	// AllErrMsg contains general error messages.
	AllErrMsg = "Error messages"
	// Error contains error messages.
	Errors = "errors"
	// ErrorNew contains new error messages.
	ErrorNew = "error"
	// ErrorCode contains error codes.
	ErrorCode = "errorCode"
	// Language contains language information.
	Language = "context-language"
	// Data contains data.
	Data = "data"
	// Success contains success message.
	Success = "Successfully added"
)

// Error messages for various scenarios.
const (
	// MongoErr contains MongoDB-related error message.
	MongoCollectionErr = "Mongo DB collection doesn't exist"

	// VerifyMessage contains verification error message.
	VerifyMessage = "Please verify whether the type, endpoint, and fields are in the correct format."

	// NotFoundErr contains error message for no matching document.
	NotFoundErr = "No matching document found."

	// QueryErr contains error message for nil query result.
	QueryErr = "Query result is nil."

	// EnLangErr contains language-related error message.
	EnLangErr = "En language not found or not of the expected language"

	// ValidationErrMessage contains validation error message.
	ValidationErrMessage = "Validation_error key not found or not of the expected type"

	// ErrorKeyErr contains error key-related error message.
	ErrorKeyErr = "Errors key not found or not of the expected type"

	// EndpointErr contains endpoint-related error message.
	EndpointErr = "Endpoint key not found or not of the expected endpoint"

	// MethodKeyErr contains method-related error message.
	MethodKeyErr = "Method key not found or not of the expected endpoint"

	// FieldErr contains field-related error message.
	FieldErr = "Fields not found or not of the expected type"

	// FailedErr contains failed updation error message.
	FailedErr = "Updation failed"

	// NoMatchErr contains no matching document error message.
	NoMatchErr = "no matching document found"

	// QueryResultErr contains query result-related error message.
	QueryResultErr = "query result is nil"

	// UrlErr contains empty URL, method, and endpoint error message.
	UrlErr = "URL, Method, and Endpoint cannot be empty"

	// NoDocumentErr contains no matching document error message.
	NoDocumentErr = "no matching document found"

	// KeyErr contains error key-related error message.
	KeyErr = "errors key not found or not of the expected type"

	// MethodErr contains method-related error message.
	MethodErr = "need method to get the field values"

	// FieldKeyErr contains field key-related error message.
	FieldKeyErr = "fields not found or not of the expected type"

	// MissingErr contains missing parameters error message.
	MissingErr = "missing parameters"

	// UnexpectedErr contains unexpected error message.
	UnexpectedErr = "unexpected error occurred"

	// FailedJsonErr contains failed to bind JSON data error message.
	FailedJsonErr = "Failed to bind JSON data"
)
