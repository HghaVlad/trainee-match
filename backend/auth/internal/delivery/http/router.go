package delivery_http

import (
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

type RouterDeps struct {
	AuthHandler *handlers.AuthHandler
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Route("/api/v1/auth", func(r chi.Router) {

		r.Post("/register", deps.AuthHandler.Register)
		r.Post("/login", deps.AuthHandler.Login)
		r.Post("/refresh", deps.AuthHandler.RefreshToken)
		r.Post("/logout", deps.AuthHandler.Logout)

	})

	return router
}
