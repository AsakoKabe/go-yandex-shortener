package handlers

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/go-chi/chi/v5"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
)

func RegisterHTTPEndpoint(router *chi.Mux, services *service.Services, cfg *config.Config) error {
	pingHandler := NewPingHandler(services.PingService)
	router.Get("/ping", pingHandler.healthDB)

	h := NewHandler(shortener.NewURLMapper(5, cfg.FileStoragePath), cfg.PrefixURL)
	router.Get("/{id}", h.getURL)
	router.Post("/", h.createShortURL)
	router.Post("/api/shorten", h.createShortURLJson)

	return nil
}
