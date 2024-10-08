package grpc

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/AsakoKabe/go-yandex-shortener/api/v1/generated"
	"github.com/AsakoKabe/go-yandex-shortener/config"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/connection"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/db/service"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"github.com/AsakoKabe/go-yandex-shortener/pkg/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// App Приложение
type App struct {
	dbPool          *sql.DB
	services        *service.Services
	handler         *Handler
	deleteWorkersWG sync.WaitGroup
}

// NewApp Конструктор для App
func NewApp(cfg *config.Config) (*App, error) {
	if cfg.DatabaseDSN == "" {
		return &App{}, nil
	}
	pool, err := connection.NewDBPool(cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Error("error to create db pool", zap.String("err", err.Error()))
		return nil, errs.ErrCreateDBPoll
	}

	pgServices, err := service.NewPostgresServices(pool)
	if err != nil {
		logger.Log.Error("error to create service", zap.String("err", err.Error()))
		return nil, errs.ErrCreateServices
	}

	return &App{
		dbPool:          pool,
		services:        pgServices,
		deleteWorkersWG: sync.WaitGroup{},
	}, nil
}

// Run Запуск приложения
func (a *App) Run(cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		slog.Error("Failed to listen: ", slog.String("err", err.Error()))
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.SlogUnaryInterceptor(slog.Default()),
			interceptors.AuthInterceptor,
		),
		grpc.ChainStreamInterceptor(
			interceptors.SlogStreamInterceptor(slog.Default()),
		),
	)

	var mapper shortener.URLShortener
	if cfg.DatabaseDSN != "" {
		mapper = shortener.NewDBUrlMapper(5, a.services.URLService)
	} else {
		mapper = shortener.NewFileURLMapper(5, cfg.FileStoragePath)
	}

	handler := NewHandler(mapper, cfg.PrefixURL, &a.deleteWorkersWG)
	a.handler = handler
	pb.RegisterShortenerServer(grpcServer, handler)

	go func() {
		slog.Info("Starting gRPC server", slog.String("addr", cfg.Addr))
		if err = grpcServer.Serve(lis); err != nil {
			slog.Error("Failed to serve: %v", slog.String("err", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			slog.Error("forcing exit")
		}
	}()

	grpcServer.GracefulStop()

	return nil
}

// Stop Завершение работы приложения
func (a *App) Stop() {
	a.handler.CloseDeleteChannel()
	a.deleteWorkersWG.Wait()
	slog.Info("delete channel closed and goroutines stopped")

	if a.dbPool == nil {
		return
	}
	err := a.dbPool.Close()
	if err != nil {
		return
	}
	slog.Info("db connection closed")
}
