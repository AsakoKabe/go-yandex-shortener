package shortener

import (
	"context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
)

type URLShortener interface {
	Add(ctx context.Context, url string, userID string) (string, error)
	AddBatch(ctx context.Context, url []string, userID string) (*[]string, error)
	Get(ctx context.Context, shortURL string) (*models.URL, bool)
	GetByUserID(ctx context.Context, userID string) (*[]models.URL, error)
	DeleteShortURLs(ctx context.Context, shortURLs []string, userID string) error
}
