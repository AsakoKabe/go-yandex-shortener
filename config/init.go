package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr            string `env:"SERVER_ADDRESS"`
	PrefixURL       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func LoadConfig() (*Config, error) {
	cfg := new(Config)

	parseFlag(cfg)

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
