package connection

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// NewDBPool функция для создания соединения с postgres
func NewDBPool(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
