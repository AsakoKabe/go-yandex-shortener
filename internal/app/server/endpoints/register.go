package endpoints

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPEndpoint(router *chi.Mux) error {
	h := NewHandler(shortener.NewURLMapper(5))

	router.Get("/{id}", h.getURL)
	router.Post("/", h.createShortURL)

	return nil
}
