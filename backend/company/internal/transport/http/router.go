package http

import (
	"log/slog"
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
	Logger         *slog.Logger
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		compmiddleware.LoggerMiddleware(deps.Logger),
	)

	router.With(compmiddleware.TimeoutMiddleware(10*time.Second)).
		Route("/api/v1/companies", func(r chi.Router) {
			// /companies/{id}
			r.With(compmiddleware.UUIDMiddleware("id")).
				Route("/{id}", func(r chi.Router) {
					r.Get("/", deps.CompanyHandler.GetByID)

					r.With(deps.AuthMiddleware.Handler,
						compmiddleware.BindJSONBodyMiddleware[dto.CompanyUpdateRequest](),
						compmiddleware.LoggingMiddleware).
						Patch("/", deps.CompanyHandler.Update)

					r.With(deps.AuthMiddleware.Handler, compmiddleware.LoggingMiddleware).
						Delete("/", deps.CompanyHandler.Delete)
				})

			r.With(deps.AuthMiddleware.Handler,
				compmiddleware.BindJSONBodyMiddleware[dto.CompanyCreateRequest](),
				compmiddleware.LoggingMiddleware).
				Post("/", deps.CompanyHandler.Create)

			r.With(compmiddleware.LoggingMiddleware).Get("/", deps.CompanyHandler.List)

			r.With(deps.AuthMiddleware.Handler, compmiddleware.LoggingMiddleware).
				Get("/me", deps.CompanyHandler.ListMy)

			// /company/{company-id}/members
			r.With(compmiddleware.UUIDMiddleware("company-id")).
				With(deps.AuthMiddleware.Handler).
				Route("/{company-id}/members", func(r chi.Router) {
					r.With(compmiddleware.BindJSONBodyMiddleware[dto.CompanyAddHrRequest](),
						compmiddleware.LoggingMiddleware).
						Post("/", deps.MemberHandler.Add)

					r.With(compmiddleware.UUIDMiddleware("user-id"),
						compmiddleware.BindJSONBodyMiddleware[dto.CompanyUpdateMemberRequest](),
						compmiddleware.LoggingMiddleware).
						Patch("/{user-id}", deps.MemberHandler.Update)

					r.With(compmiddleware.UUIDMiddleware("user-id"),
						compmiddleware.LoggingMiddleware).
						Delete("/{user-id}", deps.MemberHandler.Delete)
				})

			// /company/{company-id}/vacancies
			r.With(compmiddleware.UUIDMiddleware("company-id")).
				With(deps.AuthMiddleware.Handler).
				Route("/{company-id}/vacancies", func(r chi.Router) {
					r.With(compmiddleware.LoggingMiddleware).Get("/", deps.VacancyHandler.ListByCompany)

					r.With(compmiddleware.LoggingMiddleware).
						Get("/", deps.VacancyHandler.ListByCompany)

					r.With(compmiddleware.BindJSONBodyMiddleware[dto.VacancyCreateRequest](),
						compmiddleware.LoggingMiddleware).
						Post("/", deps.VacancyHandler.Create)

					r.With(compmiddleware.UUIDMiddleware("vacancy-id")).
						Route("/{vacancy-id}", func(r chi.Router) {
							r.With(compmiddleware.LoggingMiddleware).Get("/", deps.VacancyHandler.GetByID)

							r.With(compmiddleware.BindJSONBodyMiddleware[dto.VacancyUpdateRequest](),
								compmiddleware.LoggingMiddleware).
								Patch("/", deps.VacancyHandler.Update)

							r.With(compmiddleware.LoggingMiddleware).Post("/publish", deps.VacancyHandler.Publish)

							r.With(compmiddleware.LoggingMiddleware).Post("/archive", deps.VacancyHandler.Archive)

							r.With(compmiddleware.LoggingMiddleware).Delete("/", deps.VacancyHandler.Delete)
						})
				})
		})

	router.With(compmiddleware.TimeoutMiddleware(10*time.Second)).
		Route("/api/v1/vacancies", func(r chi.Router) {
			r.With(compmiddleware.LoggingMiddleware).Get("/", deps.VacancyHandler.List)

			r.With(compmiddleware.UUIDMiddleware("id"), compmiddleware.LoggingMiddleware).
				Get("/{id}", deps.VacancyHandler.GetPublishedByID)
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
