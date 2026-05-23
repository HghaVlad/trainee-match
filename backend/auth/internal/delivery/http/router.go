package deliveryhttp

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
)

type RouterDeps struct {
	AuthHandler  *handlers.Auth
	AdminHandler *handlers.Admin
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", deps.AuthHandler.Register)
		r.Post("/login", deps.AuthHandler.Login)
		r.Post("/refresh", deps.AuthHandler.RefreshToken)
		r.Post("/logout", deps.AuthHandler.Logout)
		r.Get("/me", deps.AuthHandler.GetMe)
	})

	router.Route("/api/v1/admin", func(r chi.Router) {
		r.Post("/new", deps.AdminHandler.NewAdmin)
	})

	return router
}
