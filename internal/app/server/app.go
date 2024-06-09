package server

import (
	"context"
	"database/sql"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/connection"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/handlers"
	"github.com/go-chi/chi/v5"
)

type App struct {
	httpServer *http.Server
	dbPool     *sql.DB
	services   *service.Services
}

func NewApp(cfg *config.Config) *App {
	pool, err := connection.NewDBPool(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return &App{
		dbPool:   pool,
		services: service.NewPostgresServices(pool),
	}
}

func (a *App) Run(cfg *config.Config) error {
	err := logger.Initialize(zap.InfoLevel)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	//router.Use(gzipMiddleware)

	err = handlers.RegisterHTTPEndpoint(router, a.services, cfg)
	if err != nil {
		log.Fatalf("Failed to register handlers: %+v", err)
		return err
	}

	a.httpServer = &http.Server{
		Addr:           cfg.Addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := http.ListenAndServe(
			cfg.Addr,
			router,
		)
		if err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)

}

func (a *App) CloseDBPool() {
	err := a.dbPool.Close()
	if err != nil {
		return
	}
}
