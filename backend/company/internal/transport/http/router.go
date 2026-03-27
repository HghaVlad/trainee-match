package http

import (
	"net/http"
	"time"

	gmiddleware "github.com/M0s1ck/g-store/src/pkg/http/middleware"
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

			extractIDFn := func(r *http.Request) string { return chi.URLParam(r, "id") }

			r.With(gmiddleware.UUIDMiddleware(extractIDFn)).
				Route("/{id}", func(r chi.Router) {

					r.Get("/", deps.CompanyHandler.GetById)

					r.With(deps.AuthMiddleware.Handler).
						With(gmiddleware.BindJSONBodyMiddleware[dto.CompanyAddHrRequest]()).
						Post("/members", deps.MemberHandler.Add)

					r.With(deps.AuthMiddleware.Handler).
						With(gmiddleware.BindJSONBodyMiddleware[dto.CompanyUpdateMemberRequest]()).
						Patch("/members/{user-id}", deps.MemberHandler.Update)

					r.With(deps.AuthMiddleware.Handler).
						Delete("/members/{user-id}", deps.MemberHandler.Delete)

					r.With(deps.AuthMiddleware.Handler).
						With(gmiddleware.BindJSONBodyMiddleware[dto.CompanyUpdateRequest]()).
						Patch("/", deps.CompanyHandler.Update)

					r.With(deps.AuthMiddleware.Handler).
						Delete("/", deps.CompanyHandler.Delete)
				})

			r.With(compmiddleware.TimeoutMiddleware(10*time.Second)).
				With(deps.AuthMiddleware.Handler).
				With(gmiddleware.BindJSONBodyMiddleware[dto.CompanyCreateRequest]()).
				Post("/", deps.CompanyHandler.Create)

			r.Get("/", deps.CompanyHandler.List)

			r.Route("/{company-id}/vacancies", func(r chi.Router) {

				r.Get("/", deps.VacancyHandler.ListByCompany)

				r.Route("/{vacancy-id}", func(r chi.Router) {

					r.With(deps.AuthMiddleware.Handler).
						Get("/", deps.VacancyHandler.GetByID)

					r.With(deps.AuthMiddleware.Handler).
						With(gmiddleware.BindJSONBodyMiddleware[dto.VacancyUpdateRequest]()).
						Patch("/", deps.VacancyHandler.Update)

					r.With(deps.AuthMiddleware.Handler).
						Post("/publish", deps.VacancyHandler.Publish)

					r.With(deps.AuthMiddleware.Handler).
						Post("/archive", deps.VacancyHandler.Archive)

					r.With(deps.AuthMiddleware.Handler).
						Delete("/", deps.VacancyHandler.Delete)
				})

				r.With(deps.AuthMiddleware.Handler).
					With(gmiddleware.BindJSONBodyMiddleware[dto.VacancyCreateRequest]()).
					Post("/", deps.VacancyHandler.Create)
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
