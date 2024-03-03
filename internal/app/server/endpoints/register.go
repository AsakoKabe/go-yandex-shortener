package endpoints

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"net/http"
)

func RegisterHTTPEndpoint(router *http.ServeMux) error {
	h := NewHandler(shortener.NewURLMapper(5))

	router.HandleFunc("/", h.root)

	return nil
}
