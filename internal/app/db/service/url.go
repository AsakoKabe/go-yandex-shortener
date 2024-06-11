package service

import (
	"context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
)

type URLService interface {
	SaveURL(ctx context.Context, url models.URL) (string, error)
	GetURL(ctx context.Context, shortURL string) (*models.URL, error)
}
