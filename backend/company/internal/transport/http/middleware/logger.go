package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	uttilslog "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
)

// LoggerMiddleware creates logger from the base and passes it into context.
func LoggerMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			innerLogger := base.With(
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", middleware.GetReqID(r.Context()),
			)

			ctx := uttilslog.WithLoggerContext(r.Context(), innerLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
