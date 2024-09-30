package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/caarlos0/env/v10"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	Addr            string `env:"SERVER_ADDRESS" json:"server_address"`
	PrefixURL       string `env:"BASE_URL" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	CertFile        string `env:"CERT_FILE"`
	KeyFile         string `env:"KEY_FILE"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	ConfigPath      string `env:"CONFIG"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// LoadConfig функция для загрузки конфигурации.
// Приоритет: переменная окружения, флаг, значение по умолчанию
func LoadConfig() (*Config, error) {
	cfg := new(Config)

	buildFlag(cfg)
	parseFlag()

	parseConfigFile(cfg)

	parseFlag() // rewrite

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}

func parseConfigFile(cfg *Config) {
	if cfg.ConfigPath == "" {
		var configPath string
		if configPath = utils.GetEnv("CONFIG", ""); configPath == "" {
			return
		}
		cfg.ConfigPath = configPath
	}

	data, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(data, cfg); err != nil {
		log.Fatal(err)
	}
}
