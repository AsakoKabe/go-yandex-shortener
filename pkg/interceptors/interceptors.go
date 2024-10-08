package interceptors

import (
	"context"
	"log/slog"
	"time"

	contextUtils "github.com/AsakoKabe/go-yandex-shortener/internal/app/context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// SlogUnaryInterceptor Custom unary interceptor for logging using slog
func SlogUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		logger.Info(
			"Unary request",
			slog.String("method", info.FullMethod),
			slog.Any("request", req),
			slog.Any("response", resp),
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

		return resp, err
	}
}

// SlogStreamInterceptor Custom stream interceptor for logging using slog
func SlogStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		logger.Info(
			"Stream request",
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.Any("error", err),
		)

		return err
	}
}

const CookieName = "jwt"

// AuthInterceptor is a gRPC interceptor for user authentication
func AuthInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	var tokenString string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	if authHeader, exists := md[CookieName]; exists && len(authHeader) > 0 {
		tokenString = authHeader[0]
	} else {
		// If no token found, generate a new one
		var err error
		tokenString, err = jwt.BuildJWTString(uuid.NewString())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate JWT token")
		}
		grpc.SendHeader(ctx, metadata.Pairs(CookieName, tokenString))
	}

	userID, err := jwt.GetUserID(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	ctx = contextUtils.SetUserID(ctx, userID)

	return handler(ctx, req)
}
