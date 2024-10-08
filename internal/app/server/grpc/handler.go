package grpc

import (
	"context"
	"errors"
	"fmt"
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
	"google.golang.org/protobuf/types/known/emptypb"
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

// CloseDeleteChannel Завершение чтения задач на удаление ссылок
func (h *Handler) CloseDeleteChannel() {
	h.delWG.Wait()
	close(h.deleteJobs)
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

func (h *Handler) AddBatchShortURL(
	ctx context.Context, req *pb.ShortenBatchRequest,
) (*pb.ShortenBatchResponse, error) {
	var originalURLs []string
	for _, batch := range req.Batch {
		originalURLs = append(originalURLs, batch.OriginalUrl)
	}

	userID := contextUtils.GetUserID(ctx)

	shortURLs, err := h.urlShortener.AddBatch(ctx, originalURLs, userID)
	if err != nil {
		logger.Log.Error("error to create short url batch", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, "failed to create short URLs")
	}

	var shortURLBatch []*pb.ShortenBatchResponse_Shorten
	for i, shortURL := range *shortURLs {
		shortURLBatch = append(
			shortURLBatch, &pb.ShortenBatchResponse_Shorten{
				ShortUrl:      h.prefixURL + shortURL,
				CorrelationId: req.Batch[i].CorrelationId,
			},
		)
	}

	return &pb.ShortenBatchResponse{
		Batch: shortURLBatch,
	}, nil
}

func (h *Handler) GetURL(ctx context.Context, req *pb.URLRequest) (*pb.URLResponse, error) {
	shortURL := req.ShortUrl

	if server.IsURLEmpty(shortURL) {
		logger.Log.Error("shortURL not found")
		return nil, status.Errorf(codes.InvalidArgument, "invalid or missing short URL")
	}

	url, ok := h.urlShortener.Get(ctx, shortURL)
	if !ok {
		logger.Log.Error("failed to retrieve URL")
		return nil, status.Errorf(codes.NotFound, "failed to retrieve URL")
	}

	if server.IsURLEmpty(url.OriginalURL) {
		logger.Log.Error("URL not found")
		return nil, status.Errorf(codes.NotFound, "original URL not found")
	}

	if url.DeletedFlag {
		return &pb.URLResponse{
			Url: "",
		}, nil
	}

	return &pb.URLResponse{
		Url: url.OriginalURL,
	}, nil

}

func (h *Handler) GetMyURLs(ctx context.Context, _ *emptypb.Empty) (
	*pb.UserURLsBatchResponse, error,
) {
	userID := contextUtils.GetUserID(ctx)
	fmt.Println(userID)

	urls, err := h.urlShortener.GetByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error("error to get URLs", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	var shortURLBatch []*pb.UserURLsBatchResponse_UserURL
	for _, url := range *urls {
		shortURLBatch = append(
			shortURLBatch, &pb.UserURLsBatchResponse_UserURL{
				ShortUrl:    h.prefixURL + url.ShortURL,
				OriginalUrl: url.OriginalURL,
			},
		)
	}

	return &pb.UserURLsBatchResponse{
		Batch: shortURLBatch,
	}, nil
}

func (h *Handler) DeleteShortURLs(ctx context.Context, req *pb.ShortURLs) (*pb.Deleted, error) {
	userID := contextUtils.GetUserID(ctx)

	h.delWG.Add(1)
	go func() {
		defer h.delWG.Done()
		h.deleteJobs <- server.DeleteJob{
			ShortURL: req.ShortUrl,
			UserID:   userID,
		}
	}()

	return &pb.Deleted{Status: true}, nil
}

func (h *Handler) GetStats(ctx context.Context, _ *emptypb.Empty) (*pb.InternalStats, error) {
	stats, err := h.urlShortener.GetStats(ctx)
	if err != nil {
		logger.Log.Error("error to get stats", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.InternalStats{
		Urls:  int32(stats.Urls),
		Users: int32(stats.Users),
	}, nil
}
