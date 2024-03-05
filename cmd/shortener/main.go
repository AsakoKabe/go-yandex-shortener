package main

import (
	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	app := server.NewApp()

	if err := app.Run(cfg); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
