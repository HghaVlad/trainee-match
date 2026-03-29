package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// LoggingMiddleware logs info about served request.
// Uses logger from context (LoggerMiddleware) with additions of previous middlewares in the chain,
// so should be called after AuthMiddleware etc.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		logger := LoggerFromContext(r.Context())

		logger.Info("http request",
			"status", ww.Status(),
			"duration", time.Since(start),
			"size", ww.BytesWritten(),
		)
	})
}
