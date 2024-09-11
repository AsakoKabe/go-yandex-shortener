package server

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/connection"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/handlers"
	middlewareUtils "github.com/AsakoKabe/go-yandex-shortener/internal/app/server/middleware"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
)

// App Приложение
type App struct {
	httpServer *http.Server
	dbPool     *sql.DB
	services   *service.Services
}

// NewApp Конструктор для App
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

// Run Запуск приложения
func (a *App) Run(cfg *config.Config) error {
	err := logger.Initialize(zap.InfoLevel)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(chiMiddleware.Logger)
	router.Use(middlewareUtils.Gzip)
	router.Use(middlewareUtils.Auth)
	router.Mount("/debug", chiMiddleware.Profiler())

	err = handlers.RegisterHTTPEndpoint(router, a.services, cfg)
	if err != nil {
		return errs.ErrRegisterEndpoints
	}

	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("cache-dir"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
	}

	a.httpServer = &http.Server{
		Addr:           cfg.Addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      manager.TLSConfig(),
	}

	go func() {
		var err error
		if cfg.EnableHTTPS {
			slog.Info("run server with HTTPS")
			err = http.ListenAndServeTLS(
				cfg.Addr,
				cfg.CertFile,
				cfg.KeyFile,
				router,
			)
		} else {
			slog.Info("run server with HTTP")
			err = http.ListenAndServe(
				cfg.Addr,
				router,
			)
		}
		if err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)

}

// CloseDBPool Закрытие соединения с БД
func (a *App) CloseDBPool() {
	if a.dbPool == nil {
		return
	}
	err := a.dbPool.Close()
	if err != nil {
		return
	}
}
