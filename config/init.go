package config

type Config struct {
	NetAddress
	PrefixURL string
}

type NetAddress struct {
	Host string
	Port string
}

func LoadConfig() (*Config, error) {
	cfg := new(Config)
	parseFlag(cfg)

	return cfg, nil
}
