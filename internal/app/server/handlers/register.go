package handlers

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/go-chi/chi/v5"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
)

func RegisterHTTPEndpoint(router *chi.Mux, services *service.Services, cfg *config.Config) error {
	var mapper shortener.URLShortener
	if cfg.DatabaseDSN != "" {
		pingHandler := NewPingHandler(services.PingService)
		router.Get("/ping", pingHandler.healthDB)
		mapper = shortener.NewDBUrlMapper(5, services.URLService)
	} else {
		mapper = shortener.NewFileURLMapper(5, cfg.FileStoragePath)
	}

	h := NewHandler(mapper, cfg.PrefixURL)
	router.Get("/{id}", h.getURL)
	router.Post("/", h.createShortURL)
	router.Post("/api/shorten", h.createShortURLJson)
	router.Post("/api/shorten/batch", h.createFromBatch)
	router.Get("/api/user/urls", h.getURLsByUser)
	router.Delete("/api/user/urls", h.deleteShorURLs)

	return nil
}
