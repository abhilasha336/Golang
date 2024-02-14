package consts

// db type,appname,accepted versions
const (
	DatabaseType     = "postgres"
	AppName          = "oauth"
	AcceptedVersions = "v1"
)

// contextacceptedversions,systemacceptedversion,contextacceptedversionindex variables
const (
	ContextAcceptedVersions       = "Accept-Version"
	ContextSystemAcceptedVersions = "System-Accept-Versions"
	ContextAcceptedVersionIndex   = "Accepted-Version-index"
)

// constant variables used in OAuth
const (
	GoogleProvider   = "google"
	FacebookProvider = "facebook"
	Provider         = "provider"
	PartnerID        = "partner_id"
	State            = "state"
	ExpTime          = 60
	TempExpTime      = 5
	RefExpTime       = 1440
	SpotifyProvider  = "spotify"
	EncryptTest      = "tuneverse-esrevenuttuneverse-tue"
)
