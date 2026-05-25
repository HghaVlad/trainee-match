package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

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

	router.Use(middleware.RequestID, middleware.RealIP)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://traineematch.space",
			"https://www.traineematch.space",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},

		ExposedHeaders: []string{
			"Link",
		},

		AllowCredentials: true,

		MaxAge: 300,
	}))

	router.Use(
		appmiddleware.LoggerMiddleware(logger),
		authMiddleware.Handler,
	)

	oapi.HandlerFromMux(
		oapi.NewStrictHandler(handler, []oapi.StrictMiddlewareFunc{appmiddleware.LoggingMiddleware}),
		router,
	)

	return router
}
