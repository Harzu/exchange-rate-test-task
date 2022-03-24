package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port     string `envconfig:"PORT"`
	LogLevel string `envconfig:"LOG_LEVEL"`
	MongoURI string `envconfig:"MONGO_URI"`
	PGDSN    string `envconfig:"PG_DSN"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
