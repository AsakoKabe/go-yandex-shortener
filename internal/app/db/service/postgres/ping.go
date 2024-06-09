package postgres

import (
	"context"
	"database/sql"
)

type PingService struct {
	db *sql.DB
}

func NewPingService(db *sql.DB) *PingService {
	return &PingService{db: db}
}

func (p *PingService) PingDB(ctx context.Context) error {
	rows, err := p.db.QueryContext(ctx, "SELECT 1;")
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()
	return err
}
