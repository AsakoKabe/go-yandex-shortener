package config

import (
	"errors"
	"flag"
	"strings"
)

func (a *NetAddress) String() string {
	return a.Host + ":" + a.Port
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	a.Host = hp[0]
	a.Port = hp[1]
	return nil
}

func parseFlag(c *Config) {
	_ = flag.Value(&c.NetAddress)
	flag.Var(&c.NetAddress, "a", "Net address host:port")
	flag.StringVar(&c.PrefixURL, "b", "http://localhost:8000", "short url prefix")
	flag.Parse()
}
