package config

import (
	"flag"
)

func parseFlag(c *Config) {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "Net address host:port")
	flag.StringVar(&c.PrefixURL, "b", "http://localhost:8080", "short url prefix")

	flag.Parse()
}
