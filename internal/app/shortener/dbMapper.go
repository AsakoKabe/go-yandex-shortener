package shortener

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	dbErrs "github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service/errs"
	handlerErrs "github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
)

type DBUrlMapper struct {
	maxLenShortURL int
	urlService     service.URLService
}

func NewDBUrlMapper(maxLenShortURL int, urlService service.URLService) *DBUrlMapper {
	return &DBUrlMapper{maxLenShortURL: maxLenShortURL, urlService: urlService}
}

func (m *DBUrlMapper) Add(ctx context.Context, originalURL string, userID string) (string, error) {
	shortURL := utils.RandStringRunes(m.maxLenShortURL)
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
	existedShortURL, err := m.urlService.SaveURL(ctx, url)
	if errors.Is(err, dbErrs.ErrOriginalURLAlreadyExist) {
		return existedShortURL, handlerErrs.ErrConflictOriginalURL
	}
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (m *DBUrlMapper) AddBatch(ctx context.Context, originalURLs []string, userID string) (*[]string, error) {
	var batchURL []models.URL
	var shortURLs []string

	for _, originalURL := range originalURLs {
		shortURL := utils.RandStringRunes(m.maxLenShortURL)
		batchURL = append(batchURL, models.URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			UserID:      userID,
		})
		shortURLs = append(shortURLs, shortURL)
	}

	err := m.urlService.SaveBatchURL(ctx, batchURL)

	if err != nil {
		return nil, err
	}

	return &shortURLs, nil
}

func (m *DBUrlMapper) Get(ctx context.Context, shortURL string) (*models.URL, bool) {
	su, err := m.urlService.GetURL(ctx, shortURL)
	if err != nil {
		logger.Log.Error("error to get short url", zap.String("err", err.Error()))
	}
	if su != nil {
		return su, true
	}

	return nil, false
}

func (m *DBUrlMapper) GetByUserID(ctx context.Context, userID string) (*[]models.URL, error) {
	return m.urlService.GetURLsByUserID(ctx, userID)
}

func (m *DBUrlMapper) DeleteShortURLs(ctx context.Context, shortURLs []string, userID string) error {
	return m.urlService.DeleteShortURLs(ctx, shortURLs, userID)
}
