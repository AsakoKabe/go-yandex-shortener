package shortener

import (
	"context"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
)

// URLShortener Интерфейс для работы со сжатыми URL
type URLShortener interface {
	Add(ctx context.Context, url string, userID string) (
		string, error,
	) // Add Сжать и добавить URL для пользователя
	AddBatch(ctx context.Context, url []string, userID string) (
		*[]string, error,
	) // AddBatch Сжать и добавить батч из URL
	Get(ctx context.Context, shortURL string) (
		*models.URL, bool,
	) // Get Получить оригинальный URL по сжатому
	GetByUserID(ctx context.Context, userID string) (
		*[]models.URL, error,
	) // GetByUserID Получить список сжатых URL по userID
	DeleteShortURLs(
		ctx context.Context, shortURLs []string, userID string,
	) error // DeleteShortURLs Удалить сжатые URl из списка
	GetStats(ctx context.Context) (
		*models.InternalStats, error,
	) // GetStats Получить статистике по сервису
}
