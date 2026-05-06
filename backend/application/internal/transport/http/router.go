package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	handler "github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
)

func NewRouter(handler *handler.Handler, logger *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	middlewares := []oapi.StrictMiddlewareFunc{}

	oapi.HandlerFromMux(
		oapi.NewStrictHandler(handler, middlewares),
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
