package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	Addr            string `env:"SERVER_ADDRESS"`
	PrefixURL       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	CertFile        string `env:"CERT_FILE"`
	KeyFile         string `env:"KEY_FILE"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
}

// LoadConfig функция для загрузки конфигурации.
// Приоритет: переменная окружения, флаг, значение по умолчанию
func LoadConfig() (*Config, error) {
	cfg := new(Config)

	parseFlag(cfg)

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
