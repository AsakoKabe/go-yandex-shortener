package middleware

import (
	"log/slog"
	"net"
	"net/http"
)

type TrustedSubnet struct {
	ipNet *net.IPNet
}

func NewTrustedSubnet(IP string) (*TrustedSubnet, error) {
	_, ipNet, err := net.ParseCIDR(IP)
	if err != nil {
		slog.Error("err to parse CIDR", slog.String("err", err.Error()))
		return nil, err
	}
	return &TrustedSubnet{ipNet: ipNet}, nil
}

// Middleware Middleware для проверки
func (t TrustedSubnet) Middleware(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			userIP := net.ParseIP(r.Header.Get("X-Real-IP"))
			if !t.ipNet.Contains(userIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}(next)
}
