package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
	"strings"

	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
)

type URLService struct {
	db *sql.DB
}

func NewURLService(db *sql.DB) (*URLService, error) {
	err := createTable(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return &URLService{db: db}, nil
}

func (u *URLService) SaveURL(ctx context.Context, url models.URL) (string, error) {
	existedURL, err := u.getURLByQuery(ctx, "select * from url WHERE original_url = $1", url.OriginalURL)
	if err != nil {
		return "", err
	}
	if existedURL != nil {
		return existedURL.ShortURL, errs.ErrOriginalURLAlreadyExist
	}

	query := `INSERT INTO url (short_url, original_url) VALUES ($1, $2)`

	_, err = u.db.ExecContext(ctx, query, url.ShortURL, url.OriginalURL)
	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}

	return "", nil
}

func (u *URLService) SaveBatchURL(ctx context.Context, batchURL []models.URL) error {
	var vals []any
	var placeholders []string
	for index, url := range batchURL {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d)",
			index*3+1,
			index*3+2))
		vals = append(vals, url.ShortURL, url.OriginalURL)
	}

	query := fmt.Sprintf("INSERT INTO url (short_url, original_url) VALUES %s", strings.Join(placeholders, ","))

	_, err := u.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil

}

func (u *URLService) GetURL(ctx context.Context, shortURL string) (*models.URL, error) {
	return u.getURLByQuery(ctx, "select * from url WHERE short_url = $1", shortURL)
}

func (u *URLService) getURLByQuery(ctx context.Context, query string, args ...any) (*models.URL, error) {
	rows, err := u.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		logger.Log.Error("error select request", zap.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	rowsExist := rows.Next()
	if !rowsExist {
		return nil, nil
	}

	var url models.URL
	if err := rows.Scan(&url.ID, &url.ShortURL, &url.OriginalURL); err != nil {
		logger.Log.Error("error parse request from db", zap.String("err", err.Error()))
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &url, nil
}

func createTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS url
	(
		id           serial primary key,
		short_url    varchar(450) NOT NULL,
		original_url varchar(450) NOT NULL UNIQUE
	)`

	_, err := db.ExecContext(ctx, query)
	return err
}
