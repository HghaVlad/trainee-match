package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/HghaVlad/trainee-match/backend/company/api/docs"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/handlers"
	compmiddleware "github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
)

type RouterDeps struct {
	CompanyHandler *handlers.CompanyHandler
	MemberHandler  *handlers.MemberHandler
	VacancyHandler *handlers.VacancyHandler
	AuthMiddleware *compmiddleware.AuthMiddleware
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	router.With(compmiddleware.TimeoutMiddleware(10*time.Second)).
		Route("/api/v1/companies", func(r chi.Router) {
			// /companies/{id}
			r.With(compmiddleware.UUIDMiddleware("id")).
				Route("/{id}", func(r chi.Router) {
					r.Get("/", deps.CompanyHandler.GetByID)

					r.With(deps.AuthMiddleware.Handler).
						With(compmiddleware.BindJSONBodyMiddleware[dto.CompanyUpdateRequest]()).
						Patch("/", deps.CompanyHandler.Update)

					r.With(deps.AuthMiddleware.Handler).
						Delete("/", deps.CompanyHandler.Delete)
				})

			r.With(deps.AuthMiddleware.Handler).
				With(compmiddleware.BindJSONBodyMiddleware[dto.CompanyCreateRequest]()).
				Post("/", deps.CompanyHandler.Create)

			r.Get("/", deps.CompanyHandler.List)

			// /company/{company-id}/members
			r.With(compmiddleware.UUIDMiddleware("company-id")).
				With(deps.AuthMiddleware.Handler).
				Route("/{company-id}/members", func(r chi.Router) {
					r.With(compmiddleware.BindJSONBodyMiddleware[dto.CompanyAddHrRequest]()).
						Post("/", deps.MemberHandler.Add)

					r.With(compmiddleware.UUIDMiddleware("user-id")).
						With(compmiddleware.BindJSONBodyMiddleware[dto.CompanyUpdateMemberRequest]()).
						Patch("/{user-id}", deps.MemberHandler.Update)

					r.With(compmiddleware.UUIDMiddleware("user-id")).
						Delete("/{user-id}", deps.MemberHandler.Delete)
				})

			// /company/{company-id}/vacancies
			r.With(compmiddleware.UUIDMiddleware("company-id")).
				With(deps.AuthMiddleware.Handler).
				Route("/{company-id}/vacancies", func(r chi.Router) {
					r.Get("/", deps.VacancyHandler.ListByCompany)

					r.With(compmiddleware.BindJSONBodyMiddleware[dto.VacancyCreateRequest]()).
						Post("/", deps.VacancyHandler.Create)

					r.With(compmiddleware.UUIDMiddleware("vacancy-id")).
						Route("/{vacancy-id}", func(r chi.Router) {
							r.Get("/", deps.VacancyHandler.GetByID)

							r.With(compmiddleware.BindJSONBodyMiddleware[dto.VacancyUpdateRequest]()).
								Patch("/", deps.VacancyHandler.Update)

							r.Post("/publish", deps.VacancyHandler.Publish)

							r.Post("/archive", deps.VacancyHandler.Archive)

							r.Delete("/", deps.VacancyHandler.Delete)
						})
				})
		})

	router.With(compmiddleware.TimeoutMiddleware(10*time.Second)).
		Route("/api/v1/vacancies", func(r chi.Router) {
			r.Get("/", deps.VacancyHandler.List)

			r.Get("/{vacancy-id}", deps.VacancyHandler.GetPublishedByID)
		})

	addHello(router)

	addSwagger(router)

	return router
}

func addHello(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
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
