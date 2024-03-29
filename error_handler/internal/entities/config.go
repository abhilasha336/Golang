package entities

// EnvConfig represents the configuration for the environment.
type EnvConfig struct {
	Debug            bool         `default:"true" split_words:"true"`
	Port             int          `default:"8080" split_words:"true"`
	Db               Database     `split_words:"true"`
	PsqlDb           PsqlDatabase `split_words:"true"`
	AcceptedVersions []string     `required:"true" split_words:"true"`
	MigrationPath    string       ` split_words:"true"`
	LanguageUrl      string       ` split_words:"true"`
	DocID            string       ` split_words:"true"`
	EndpointDocID    string       ` split_words:"true"`
}

// Database represents the database configuration.
type Database struct {
	Driver    string `default:"mongodb" split_words:"true"`
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string
	MaxActive int
	MaxIdle   int
}

// Psql database connection struct
type PsqlDatabase struct {
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string
	MaxActive int
	MaxIdle   int
}
