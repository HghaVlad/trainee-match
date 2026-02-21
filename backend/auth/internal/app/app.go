package app

import (
	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
	deliveryhttp "github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/keycloack"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/services"
	"net/http"
)

type App struct {
	Config     *config.Config
	httpRouter http.Handler
}

func Build(conf *config.Config) *App {

	keycloackClient := keycloack.NewClient(
		conf.KeyCloack.URL, conf.KeyCloack.Realm, conf.KeyCloack.ClientID, conf.KeyCloack.ClientSecret, conf.KeyCloack.AdminUsername, conf.KeyCloack.AdminPassword,
	)

	authService := services.NewAuth(keycloackClient)

	deps := deliveryhttp.RouterDeps{
		AuthHandler: handlers.NewAuthHandler(authService, conf.KeyCloack.AccessTokenExpires, conf.KeyCloack.RefreshTokenExpires),
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
