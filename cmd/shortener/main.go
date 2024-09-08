package main

import (
	"fmt"
	"log"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	printBuildInfo()

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

func printBuildInfo() {
	fmt.Println("Build version: " + buildVersion)
	fmt.Println("Build date: " + buildDate)
	fmt.Println("Build commit: " + buildCommit)
}
