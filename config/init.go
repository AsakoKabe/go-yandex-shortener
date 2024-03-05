package config

type Config struct {
	NetAddress
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
