package service

import (
	"context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
)

// URLService интерфейс для работы с сервисом. Реализуется разными БД
type URLService interface {
	SaveURL(ctx context.Context, url models.URL) (string, error) // SaveURL Сохранить URL в БД
	SaveBatchURL(
		ctx context.Context, batchURL []models.URL,
	) error // SaveBatchURL Сохранить батч из URL в БД
	GetURL(ctx context.Context, shortURL string) (
		*models.URL, error,
	) // GetURL Получить исходный URL по сжатому
	GetURLsByUserID(ctx context.Context, userID string) (
		*[]models.URL, error,
	) // GetURLsByUserID Получить список всех URL для пользователя
	DeleteShortURLs(
		ctx context.Context, shortURLs []string, userID string,
	) error // DeleteShortURLs Удалить список из URL для пользователя
}
