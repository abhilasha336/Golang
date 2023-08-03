package config

import (
	"backend-code-base-template/internal/entities"

	"github.com/kelseyhightower/envconfig"
)

// LoadConfig
func LoadConfig(appName string) (*entities.EnvConfig, error) {
	var cfg entities.EnvConfig
	err := envconfig.Process(appName, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
