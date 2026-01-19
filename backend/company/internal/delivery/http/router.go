package delivery_http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"

	"github.com/HghaVlad/trainee-match/backend/company/api/docs"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
)

type RouterDeps struct {
	ProfileHandler *handlers.CompanyHandler
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	router.Route("/api/v1/companies", func(r chi.Router) {

		r.With(my_middleware.UUIDMiddleware("id")).
			Route("/{id}", func(r chi.Router) {

				// TODO: maybe later change to /profile
				r.Get("/", deps.ProfileHandler.GetById)
			})

	})

	addHello(router)
	addSwagger(router)

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

func addSwagger(r chi.Router) {
	r.Get("/swagger/*", func(w http.ResponseWriter, req *http.Request) {
		// Подставляем host динамически
		docs.SwaggerInfo.Host = req.Host
		docs.SwaggerInfo.Schemes = []string{"http"}

		// Учитываем gateway prefix
		if p := req.Header.Get("X-Forwarded-Prefix"); p != "" {
			docs.SwaggerInfo.BasePath = p + "/api/v1"
		} else {
			docs.SwaggerInfo.BasePath = "/api/v1"
		}

		httpSwagger.WrapHandler(w, req)
	})
}
