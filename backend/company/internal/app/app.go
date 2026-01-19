package app

import (
	"log"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/get_company"
)

type App struct {
	httpRouter http.Handler
	conf       *config.Config
}

func Build(conf *config.Config) *App {
	psgConf := infra_postgres.NewConfig(conf)
	logger := service_logger.NewSlogLogger()
	compDB, err := infra_postgres.New(psgConf, logger)
	if err != nil {
		log.Fatal(err)
	}

	compRepo := repository.NewCompanyRepository(compDB)

	profileGetByIDUc := get_company.NewGetByIDUsecase(compRepo)

	profileHandler := handlers.NewProfileHandler(profileGetByIDUc)

	routerDeps := &delivery_http.RouterDeps{
		ProfileHandler: profileHandler,
	}

	httpRouter := delivery_http.NewRouter(routerDeps)

	return &App{
		httpRouter: httpRouter,
		conf:       conf,
	}
}

func (app *App) Run() {
	err := http.ListenAndServe(app.conf.HTTP.Addr, app.httpRouter)
	if err != nil {
		panic(err)
	}
}
