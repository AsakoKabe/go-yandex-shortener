package main

import (
	"log"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer app.CloseDBPool()

	if err := app.Run(cfg); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
