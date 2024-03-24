package logger

import (
	"github.com/go-chi/httplog/v2"
	"log/slog"
)

var Logger = httplog.NewLogger(
	"log-shortener",
	httplog.Options{
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		MessageFieldName: "message",
	},
)

func Initialize(level slog.Level) {
	Logger = httplog.NewLogger("log-shortener", httplog.Options{
		LogLevel:         level,
		Concise:          true,
		MessageFieldName: "message",
	})
}
