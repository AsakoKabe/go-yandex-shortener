package grpc

import (
	"context"
	"errors"
	"sync"

	pb "github.com/AsakoKabe/go-yandex-shortener/api/v1/generated"
	contextUtils "github.com/AsakoKabe/go-yandex-shortener/internal/app/context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedShortenerServer
	urlShortener shortener.URLShortener
	prefixURL    string
	deleteJobs   chan server.DeleteJob
	delWG        *sync.WaitGroup
}

func NewHandler(
	urlShortener shortener.URLShortener,
	prefixURL string,
	deleteWorkersWG *sync.WaitGroup,
) *Handler {
	jobs := make(chan server.DeleteJob, server.NumDeleteJobs)

	for w := 1; w <= server.NumWorkers; w++ {
		deleteWorkersWG.Add(1)
		go server.DeleteWorker(deleteWorkersWG, urlShortener, jobs)
	}

	return &Handler{
		urlShortener: urlShortener,
		prefixURL:    prefixURL,
		deleteJobs:   jobs,
		delWG:        &sync.WaitGroup{},
	}
}

func (h *Handler) AddShortURL(ctx context.Context, req *pb.ShortenRequest) (
	*pb.ShortenResponse, error,
) {
	originalURL := req.Url

	if originalURL == "" {
		return nil, status.Errorf(codes.InvalidArgument, "original URL is empty")
	}

	userID := contextUtils.GetUserID(ctx)

	shortURL, err := h.urlShortener.Add(ctx, originalURL, userID)
	if errors.Is(err, errs.ErrConflictOriginalURL) {
		logger.Log.Info("original url already exists", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.AlreadyExists, "original URL already exists")
	} else if err != nil {
		logger.Log.Error("error creating short URL", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, "failed to create short URL")
	}

	return &pb.ShortenResponse{
		ShortUrl: h.prefixURL + shortURL,
	}, nil
}

// func (h *Handler) AddBatchShortURL(
// 	context.Context, *pb.ShortenBatchRequest,
// ) (*pb.ShortenBatchResponse, error) {
// }
//
// func (h *Handler) GetURL(context.Context, *pb.URLRequest) (*pb.URLResponse, error) {
// }
//
// func (h *Handler) GetMyURLs(context.Context, *emptypb.Empty) (*pb.UserURLsBatchResponse, error) {
// }
//
// func (h *Handler) DeleteShortURLs(context.Context, *pb.ShortURLs) (*pb.Deleted, error) {
// }
//
// func (h *Handler) GetStats(context.Context, *emptypb.Empty) (*pb.InternalStats, error) {
// }

// CloseDeleteChannel Завершение чтения задач на удаление ссылок
func (h *Handler) CloseDeleteChannel() {
	h.delWG.Wait()
	close(h.deleteJobs)
}
