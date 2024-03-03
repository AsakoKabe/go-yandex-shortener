package main

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server"
	"log"
)

func main() {
	app := server.NewApp()
	port := "8080"

	if err := app.Run(port); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
