package handlers

import (
	"log/slog"
	"sync"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
)

// RegisterHTTPEndpoint Функция для регистрации endpoints
func RegisterHTTPEndpoint(
	deleteWG *sync.WaitGroup, router *chi.Mux, services *service.Services, cfg *config.Config,
) (*Handler, error) {
	var mapper shortener.URLShortener
	if cfg.DatabaseDSN != "" {
		pingHandler := NewPingHandler(services.PingService)
		router.Get("/ping", pingHandler.healthDB)
		mapper = shortener.NewDBUrlMapper(5, services.URLService)
	} else {
		mapper = shortener.NewFileURLMapper(5, cfg.FileStoragePath)
	}

	trustedSubnet, err := middleware.NewTrustedSubnet(cfg.TrustedSubnet)
	if err != nil {
		slog.Error("error to create trusted subnet", slog.String("err", err.Error()))
		return nil, err
	}

	h := NewHandler(deleteWG, mapper, cfg.PrefixURL)
	router.Get("/{id}", h.getURL)
	router.Post("/", h.createShortURL)
	router.Post("/api/shorten", h.createShortURLJson)
	router.Post("/api/shorten/batch", h.createFromBatch)
	router.Get("/api/user/urls", h.getURLsByUser)
	router.Delete("/api/user/urls", h.deleteShorURLs)
	router.Route(
		"/api/internal/stats", func(r chi.Router) {
			r.Use(trustedSubnet.Middleware)
			r.Get("/", h.getStats)
		},
	)

	return h, nil
}
