package service

import (
	"database/sql"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/postgres"
)

type Services struct {
	PingService PingService
	URLService  URLService
}

func NewPostgresServices(db *sql.DB) *Services {
	return &Services{
		PingService: postgres.NewPingService(db),
		URLService:  postgres.NewURLService(db),
	}
}
