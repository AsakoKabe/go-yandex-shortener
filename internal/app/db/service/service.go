package service

import (
	"database/sql"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/postgres"
)

type Services struct {
	PingService PingService
}

func NewPostgresServices(db *sql.DB) *Services {
	return &Services{
		PingService: postgres.NewPingService(db),
	}
}
