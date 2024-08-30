package postgres

import (
	"context"
	"database/sql"
)

// PingService структура реализации сервиса к БД для postgres
type PingService struct {
	db *sql.DB
}

// NewPingService конструктор для PingService
func NewPingService(db *sql.DB) *PingService {
	return &PingService{db: db}
}

// PingDB функция для отправки пинга в БД для postgres
func (p *PingService) PingDB(ctx context.Context) error {
	return p.db.PingContext(ctx)
}
