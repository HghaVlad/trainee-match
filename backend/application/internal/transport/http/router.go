package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/handlers"
	appmiddleware "github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
)

func NewRouter(
	handler *handlers.Handler,
	authMiddleware *appmiddleware.AuthMiddleware,
	logger *slog.Logger,
) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		appmiddleware.LoggerMiddleware(logger),
		authMiddleware.FakeHandler, // TODO: add real handler
	)

	oapi.HandlerFromMux(
		oapi.NewStrictHandler(handler, []oapi.StrictMiddlewareFunc{appmiddleware.LoggingMiddleware}),
		router,
	)

	router.Route("/api/v1/application", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// Temporary handler for testing
			r.Get("/test", func(w http.ResponseWriter, _ *http.Request) {
				_, err := w.Write([]byte("Hello World"))
				if err != nil {
					return
				}
			})
		})
	})

	return router
}
