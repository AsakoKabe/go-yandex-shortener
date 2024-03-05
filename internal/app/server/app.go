package server

import (
	"context"
	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/endpoints"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(cfg *config.Config) error {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	err := endpoints.RegisterHTTPEndpoint(router)
	if err != nil {
		log.Fatalf("Failed to register endpoints: %+v", err)
		return err
	}

	a.httpServer = &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := http.ListenAndServe(
			cfg.Host+":"+cfg.Port,
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
