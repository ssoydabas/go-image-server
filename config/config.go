package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"production"`
	Server      struct {
		Port            string        `env:"SERVER_PORT" envDefault:"8080"`
		MaxFileSize     int64         `env:"MAX_FILE_SIZE" envDefault:"10485760"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	}
	Storage struct {
		BasePath string `env:"STORAGE_PATH" envDefault:"./data"`
		DevPath  string `env:"DEV_STORAGE_PATH" envDefault:"./dev-data"`
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
