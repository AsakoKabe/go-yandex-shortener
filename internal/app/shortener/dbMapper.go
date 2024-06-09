package shortener

import (
	"context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
)

type DBUrlMapper struct {
	maxLenShortURL int
	urlService     service.URLService
}

func NewDBUrlMapper(maxLenShortURL int, urlService service.URLService) *DBUrlMapper {
	return &DBUrlMapper{maxLenShortURL: maxLenShortURL, urlService: urlService}
}

func (m *DBUrlMapper) Add(ctx context.Context, originalURL string) (string, error) {
	shortURL := "/" + utils.RandStringRunes(m.maxLenShortURL)
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	err := m.urlService.SaveURL(ctx, url)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (m *DBUrlMapper) Get(ctx context.Context, shortURL string) (string, bool) {
	su, err := m.urlService.GetURL(ctx, shortURL)
	if err != nil {
		logger.Log.Error("error to create short url", zap.String("err", err.Error()))
	}
	if su != nil {
		return su.OriginalURL, true
	}
	return "", false
}
