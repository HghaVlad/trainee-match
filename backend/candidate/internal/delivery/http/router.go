package http

import (
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type RouterDeps struct {
	authMiddleware   *auth.Middleware
	candidateHandler *handlers.Candidate
}

func NewRouterDeps(authMiddleware *auth.Middleware, candidateHandler *handlers.Candidate) *RouterDeps {
	return &RouterDeps{
		authMiddleware:   authMiddleware,
		candidateHandler: candidateHandler,
	}
}

func NewRouter(deps *RouterDeps) http.Handler {

	router := chi.NewRouter()

	router.Route("/api/v1/candidate", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(deps.authMiddleware.Handler)
			r.Get("/me", deps.candidateHandler.GetMe)
			r.Post("/", deps.candidateHandler.CreateCandidate)
			r.Patch("/", deps.candidateHandler.UpdateCandidate)
		})
	})

	router.Get("/swagger/*", handlers.SwaggerHandler)

	return router
}
