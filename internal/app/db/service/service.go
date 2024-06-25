package service

import (
	"database/sql"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/postgres"
)

type Services struct {
	PingService PingService
	URLService  URLService
}

func NewPostgresServices(db *sql.DB) (*Services, error) {
	urlService, err := postgres.NewURLService(db)
	if err != nil {
		return nil, err
	}
	return &Services{
		PingService: postgres.NewPingService(db),
		URLService:  urlService,
	}, nil
}
