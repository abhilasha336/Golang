package consts

const (
	// AcceptedVersions is the accepted version string.
	AcceptedVersions = "v1.0"

	ContextAcceptedVersions       = "Accept-Version"
	ContextSystemAcceptedVersions = "System-Accept-Versions"
	ContextAcceptedVersionIndex   = "Accepted-Version-index"
	// Context keys for locale and language.
	ContextLocaleLang = "lan"
	HeaderLanguage    = "Accept-Language"
	ContextEndPoints  = "context-endpoints"

	ContextErrorResponses       = "context-error-response"
	HeaderLocallizationLanguage = "Accept-Language"
)

var (
	CacheErrorData    = "CACHE_ERROR_DATA"
	CacheEndPointData = "CACHE_ENDPOINT_DATA"
)

// Context setting values
const (
	// ContextErrorResponses        = "context-error-response"
	ContextLocallizationLanguage = "lan"
)

// KeyNames
const (
	ValidationErr     = "validation_error"
	ForbiddenErr      = "forbidden"
	UnauthorisedErr   = "unauthorized"
	NotFound          = "not found"
	InternalServerErr = "internal_server_error"
	Errors            = "errors"
	AllError          = "AllError"
	Registration      = "registration"
	ErrorCode         = "errorCode"
	MemberIDErr       = "member_id"
	Message           = "message"
	Language          = "context-language"
)

const (
	EndpointErr = "Error occured while loading endpoints from service"
	ContextErr  = "Error occured while loading error from service"
)

// Default values
const (
	DefaultDateFormat = "02-01-2006"
	DefaultDirName    = "logs"
	DefaultPage       = 1
	DefaultLimit      = 10
	MaxLimit          = 100
	MaxLength         = 20
	MinLength         = 3
)

type contextKey string

const (
	LogData contextKey = "log_data"
)
const (
	ContextRequestID          = "req_id"
	ContextRequestURI         = "uri"
	ContextRequestMethod      = "method"
	ContextRequestIP          = "user_ip"
	ContextRequestStatus      = "response_code"
	ContextRequestTimetaken   = "duration_ms"
	ContextService            = "service"
	ContextRequestURITemplate = "endpoint"
	ContextResponseDump       = "response_dump"
	ContextRequestDump        = "request_dump"
	ContextTimeStamp          = "timestamp"
	ContextLogLevel           = "log_level"
	ContextMessage            = "message"
)
const (
	Email             = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	MaxURLRuneCount   = 2083
	MinURLRuneCount   = 3
	URLPath           = `((\/|\?|#)[^\s]*)`
	IP                = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLPort           = `(:(\d{1,5}))`
	URLSchema         = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername       = `(\S+(:\S*)?@)`
	URLIP             = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	URLSubdomain      = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URL               = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	ValidTitlePattern = "^[a-zA-Z0-9_-]+$"
)

// Constants used for regex validation
type Validation struct {
	SpecialCharKey    string
	URLKey            string
	DateKey           string
	EmailKey          string
	JsonPathKey       string
	TrackExtensionKey string
	Isrc              string
	Iswc              string
	ReplaceString     string
	Uppercase         string
	LyricsEnd         string
}

var Validations = Validation{
	SpecialCharKey:    "specialchar",
	URLKey:            "url",
	DateKey:           "date",
	EmailKey:          "email",
	JsonPathKey:       "https://tuneversev2.s3.amazonaws.com/RegexPattern/pattern.json",
	TrackExtensionKey: "trackextension",
	Isrc:              "isrc",
	Iswc:              "iswc",
	ReplaceString:     "replace",
	Uppercase:         "uppercase",
	LyricsEnd:         "lyricsend",
}

const LogUrl = "activitylog"

// Time constants
const (
	InputFormat  = "2006-01-02T15:04:05.999999Z"
	OutputFormat = "02 January 2006 03:04 PM"
)

var LengthConstraints = map[string]string{
	"max": "max_length",
	"min": "min_length",
}
