package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/HghaVlad/trainee-match/backend/application/internal/delivery/http/handlers"
)

type RouterDeps struct {
	ApplicationHandler handlers.Application
}

func NewRouterDeps() *RouterDeps {
	return &RouterDeps{}
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

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
