package entities

type EnvConfig struct {
	Debug            bool     `default:"true" split_words:"true"`
	Port             int      `default:"3000" split_words:"true"`
	Db               Database `split_words:"true"`
	AcceptedVersions []string `required:"true" split_words:"true"`
}

type Database struct {
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string `envconfig:"default=public"`
	MaxActive int
	MaxIdle   int
}
