package main

import (
	"fmt"
	"log"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}
	fmt.Println(cfg.DatabaseDSN)
	app := server.NewApp(cfg)
	defer app.CloseDBPool()

	if err := app.Run(cfg); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
