package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type loggerCtxKeyT struct{}

//nolint:gochecknoglobals // ctx key
var loggerCtxKey = loggerCtxKeyT{}

// LoggerMiddleware creates logger from the base and passes it into context.
func LoggerMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			innerLogger := base.With(
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", middleware.GetReqID(r.Context()),
			)

			ctx := context.WithValue(r.Context(), loggerCtxKey, innerLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtxKey).(*slog.Logger)
	if !ok {
		panic("logger not found in context. logger middleware is not applied")
	}

	return logger
}
