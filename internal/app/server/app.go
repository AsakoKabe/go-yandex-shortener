package server

import (
	"context"
	"database/sql"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/connection"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/handlers"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
)

type App struct {
	httpServer *http.Server
	dbPool     *sql.DB
	services   *service.Services
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg.DatabaseDSN == "" {
		return &App{}, nil
	}
	pool, err := connection.NewDBPool(cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Error("error to create db pool", zap.String("err", err.Error()))
		return nil, errs.ErrCreateDBPoll
	}

	pgServices, err := service.NewPostgresServices(pool)
	if err != nil {
		logger.Log.Error("error to create service", zap.String("err", err.Error()))
		return nil, errs.ErrCreateServices
	}

	return &App{
		dbPool:   pool,
		services: pgServices,
	}, nil
}

func (a *App) Run(cfg *config.Config) error {
	err := logger.Initialize(zap.InfoLevel)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(gzipMiddleware)

	err = handlers.RegisterHTTPEndpoint(router, a.services, cfg)
	if err != nil {
		return errs.ErrRegisterEndpoints
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
	if a.dbPool == nil {
		return
	}
	err := a.dbPool.Close()
	if err != nil {
		return
	}
}
