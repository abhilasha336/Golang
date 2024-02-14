package entities

// EnvConfig represents the configuration structure for the application.
type EnvConfig struct {
	Debug                  bool     `default:"true" split_words:"true"`  // Flag indicating debug mode (default: true)
	Port                   int      `default:"8039" split_words:"true"`  // Port for server to listen on (default: 8080)
	Db                     Database `split_words:"true"`                 // Database configuration
	AcceptedVersions       []string `required:"true" split_words:"true"` // List of accepted API versions (required)
	MigrationPath          string   `split_words:"true"`                 // Path to migration files
	LocalisationServiceURL string   `split_words:"true"`                 // URL of the localization service
	EndpointURL            string   `split_words:"true"`                 // URL of the endpoint service
	LoggerServiceURL       string   `envconfig:"LOGGER_SERVICE_URL"`
	LoggerSecret           string   `envconfig:"LOGGER_SECRET"`
	JwtKey                 string   `split_words:"true"`
	DecryptionKey          string   `split_words:"true"`
}

// Database represents the configuration for the database connection.
type Database struct {
	User      string // Database username
	Password  string // Database password
	Port      int    // Database port
	Host      string // Database host
	DATABASE  string // Database name
	Schema    string // Database schema
	MaxActive int    // Maximum number of active connections
	MaxIdle   int    // Maximum number of idle connections
}
