package http

import (
	"log/slog"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/handlers"
	appmiddleware "github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	handler *handlers.Handler,
	authMiddleware *appmiddleware.AuthMiddleware,
	logger *slog.Logger,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID, middleware.RealIP)

	router.Use(appmiddleware.Cors())

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	router.Group(func(r chi.Router) {
		r.Use(
			appmiddleware.LoggerMiddleware(logger),
			authMiddleware.Handler,
		)

		oapi.HandlerFromMux(
			oapi.NewStrictHandler(handler, []oapi.StrictMiddlewareFunc{appmiddleware.LoggingMiddleware}),
			r,
		)
	})

	return router
}
