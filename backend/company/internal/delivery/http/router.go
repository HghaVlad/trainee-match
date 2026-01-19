package delivery_http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
)

type RouterDeps struct {
	ProfileHandler *handlers.ProfileHandler
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	router.Route("/api/v1/company/profile", func(r chi.Router) {

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", deps.ProfileHandler.GetById)
		})

	})

	addHello(router)

	return router
}

func addHello(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			return
		}
	})
}
