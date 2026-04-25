package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	utilslog "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
)

// LoggingMiddleware logs info about served request.
// Uses logger from context (LoggerMiddleware) with additions of previous middlewares in the chain,
// so should be called after AuthMiddleware etc.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		ctx := r.Context()
		logger := utilslog.FromContext(ctx)

		duration := time.Since(start)

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Error("request timeout",
				"status", ww.Status(),
				"duration", duration,
			)
			return
		}

		logger.Info("http request",
			"status", ww.Status(),
			"duration", duration,
			"size", ww.BytesWritten(),
		)
	})
}
