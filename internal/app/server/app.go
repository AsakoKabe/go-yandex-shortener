package server

import (
	"context"
	"log/slog"
	"sync"

	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
)

type App interface {
	Run(cfg *config.Config) error
	Stop()
}

const NumDeleteJobs = 5
const NumWorkers = 5

type DeleteJob struct {
	ShortURL []string
	UserID   string
}

func IsURLEmpty(url string) bool {
	return url == ""
}

func DeleteWorker(
	wg *sync.WaitGroup, urlShortener shortener.URLShortener, jobs <-chan DeleteJob,
) {
	defer wg.Done()
	for j := range jobs {
		err := urlShortener.DeleteShortURLs(context.Background(), j.ShortURL, j.UserID)
		if err != nil {
			logger.Log.Error("error to delete url", zap.String("err", err.Error()))
		}
	}
	slog.Info("stop delete worker")
}
