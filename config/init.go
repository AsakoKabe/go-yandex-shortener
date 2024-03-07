package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

type Config struct {
	Addr      string `env:"SERVER_ADDRESS"`
	PrefixURL string `env:"BASE_URL"`
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
