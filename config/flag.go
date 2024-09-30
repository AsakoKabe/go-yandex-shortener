package config

import (
	"flag"
)

func buildFlag(c *Config) {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "Net address host:port")
	flag.StringVar(&c.PrefixURL, "b", "http://localhost:8080", "short url prefix")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&c.DatabaseDSN, "d", "", "db path")
	flag.BoolVar(&c.EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&c.ConfigPath, "c", "", "path to json config")
	flag.StringVar(&c.TrustedSubnet, "t", "", "trusted subnet")
}

func parseFlag() {
	flag.Parse()
}
