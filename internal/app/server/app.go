package server

import (
	"context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/endpoints"
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

func (a *App) Run(port string) error {
	router := http.NewServeMux()

	err := endpoints.RegisterHTTPEndpoint(router)
	if err != nil {
		log.Fatalf("Failed to register endpoints: %+v", err)
		return err
	}

	a.httpServer = &http.Server{
		Addr:           "localhost:" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := http.ListenAndServe(
			":"+port,
			Conveyor(router, logRequest),
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
