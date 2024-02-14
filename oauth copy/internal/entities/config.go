package entities

// EnvConfig struct used to store env variables from .env file
type EnvConfig struct {
	Debug            bool     `default:"true" split_words:"true"`
	Port             int      `default:"8080" split_words:"true"`
	Db               Database `split_words:"true"`
	AcceptedVersions []string `required:"true" split_words:"true"`
	MigrationPath    string   `split_words:"true"`
	EncryptionKey    string   `split_words:"true"`
	JwtKey           string   `split_words:"true"`
	LoggerServiceURL string   `split_words:"true"`
	LoggerSecret     string   `split_words:"true"`
	MemberServiceURL string   `split_words:"true"`
	ServiceURL       string   `split_words:"true"`
}

// Database struct used to store db's env variables from .env file
type Database struct {
	User     string
	Password string
	Port     int
	Host     string
	DATABASE string
	Schema   string
	// Schema    string `envconfig:"default=public"`
	MaxActive int
	MaxIdle   int
}
