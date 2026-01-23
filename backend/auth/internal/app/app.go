package app

import (
	"github.com/HghaVlad/trainee-match/backend/auth/internal/auth"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
	deliveryhttp "github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
	"net/http"
)

type App struct {
	Config     *config.Config
	httpRouter http.Handler
}

func Build(conf *config.Config) *App {

	keycloackClient := auth.NewKeycloakClient(
		conf.KeyCloack.URL, conf.KeyCloack.Realm, conf.KeyCloack.ClientID, conf.KeyCloack.ClientSecret, conf.KeyCloack.AdminUsername, conf.KeyCloack.AdminPassword,
	)

	deps := deliveryhttp.RouterDeps{
		AuthHandler: handlers.NewAuthHandler(keycloackClient),
	}
	httpRouter := deliveryhttp.NewRouter(&deps)

	return &App{conf, httpRouter}
}

func (app *App) Run() {
	err := http.ListenAndServe(app.Config.Addr, app.httpRouter)
	if err != nil {
		panic(err)
	}
}
