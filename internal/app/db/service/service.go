package service

import (
	"database/sql"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/postgres"
)

// Services Набор сервисов для доступов к БД
type Services struct {
	PingService PingService
	URLService  URLService
}

// NewPostgresServices функция для создания сервисов postgres
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
