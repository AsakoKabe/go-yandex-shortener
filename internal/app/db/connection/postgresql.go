package connection

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDBPool(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	return pool, nil

}
