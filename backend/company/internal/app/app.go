package app

import (
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
)

type App struct {
	httpRouter http.Handler
	conf       *config.Config
}

func Build(conf *config.Config) *App {
	psgConf := infra_postgres.NewConfig(conf)
	_ = psgConf

	profileHandler := handlers.NewProfileHandler()

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
