package delivery_http

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
)

type RouterDeps struct {
	AuthHandler *handlers.Auth
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Route("/api/v1/auth", func(r chi.Router) {

		r.Post("/register", deps.AuthHandler.Register)
		r.Post("/login", deps.AuthHandler.Login)
		r.Post("/refresh", deps.AuthHandler.RefreshToken)
		r.Post("/logout", deps.AuthHandler.Logout)
		r.Post("/me", deps.AuthHandler.GetMe)

	})

	return router
}
