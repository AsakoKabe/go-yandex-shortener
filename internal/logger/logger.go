package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strconv"
	"time"
)

var Log = zap.NewNop()

func Initialize(level zapcore.Level) error {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				Log.Info(
					"request",
					zap.String("URI", r.RequestURI),
					zap.String("method", r.Method),
					zap.String("elapsed", strconv.FormatInt(int64(time.Since(t1)), 10)),
				)
				Log.Info(
					"response",
					zap.String("status", strconv.Itoa(ww.Status())),
					zap.String("bytes", strconv.Itoa(ww.BytesWritten())),
				)
				//entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}(next)
}
