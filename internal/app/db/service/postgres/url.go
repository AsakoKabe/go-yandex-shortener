package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
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
	existedURLs, err := u.getURLsByQuery(ctx, "select * from url WHERE original_url = $1", url.OriginalURL)
	if err != nil {
		return "", err
	}
	existedURL := getFirstURL(existedURLs)
	if existedURL != nil {
		return existedURL.ShortURL, errs.ErrOriginalURLAlreadyExist
	}

	query := `INSERT INTO url (user_id, short_url, original_url) VALUES ($1, $2, $3)`
	_, err = u.db.ExecContext(ctx, query, url.UserID, url.ShortURL, url.OriginalURL)
	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}

	return "", nil
}

func (u *URLService) SaveBatchURL(ctx context.Context, batchURL []models.URL) error {
	var vals []any
	var placeholders []string
	for index, url := range batchURL {
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d,$%d)",
			index*3+1, index*3+2, index*3+3))
		vals = append(vals, url.UserID, url.ShortURL, url.OriginalURL)
	}

	query := fmt.Sprintf("INSERT INTO url (user_id, short_url, original_url) VALUES %s", strings.Join(placeholders, ","))

	_, err := u.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil

}

func (u *URLService) GetURL(ctx context.Context, shortURL string) (*models.URL, error) {
	urls, err := u.getURLsByQuery(ctx, "select * from url WHERE short_url = $1", shortURL)
	if err != nil {
		return nil, err
	}

	return getFirstURL(urls), nil
}

func (u *URLService) getURLsByQuery(ctx context.Context, query string, args ...any) (*[]models.URL, error) {
	rows, err := u.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		logger.Log.Error("error select urls", zap.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var urls []models.URL

	for rows.Next() {
		err = u.parseRow(rows, &urls)
		if err != nil {
			return &urls, err
		}
	}

	if err = rows.Err(); err != nil {
		return &urls, err
	}
	return &urls, nil
}

func (u *URLService) parseRow(rows *sql.Rows, nameTasks *[]models.URL) error {
	var url models.URL
	if err := rows.Scan(&url.ID, &url.UserID, &url.ShortURL, &url.OriginalURL); err != nil {
		logger.Log.Error("error parse request from db", zap.String("err", err.Error()))
		return err
	}
	*nameTasks = append(*nameTasks, url)

	return nil
}

func getFirstURL(urls *[]models.URL) *models.URL {
	if len(*urls) == 0 {
		return nil
	}

	return &(*urls)[0]
}

func (u *URLService) GetURLsByUserID(ctx context.Context, userID string) (*[]models.URL, error) {
	return u.getURLsByQuery(ctx, "select * from url WHERE user_id = $1", userID)
}

func createTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS url
	(
		id           serial primary key,
		user_id 	 uuid,
		short_url    varchar(450) NOT NULL,
		original_url varchar(450) NOT NULL UNIQUE
	)`

	_, err := db.ExecContext(ctx, query)
	return err
}
