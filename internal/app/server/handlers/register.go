package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
)

func RegisterHTTPEndpoint(router *chi.Mux, cfg *config.Config) error {
	h := NewHandler(shortener.NewURLMapper(5), cfg.PrefixURL)

	router.Get("/{id}", h.getURL)
	router.Post("/", h.createShortURL)
	router.Post("/api/shorten", h.createShortURLJson)

	return nil
}
