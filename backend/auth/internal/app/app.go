package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
	deliveryhttp "github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/keycloack"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/message_broker/kafka"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/message_broker/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/services"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/services/outbox"
)

type App struct {
	Config      *config.Config
	httpRouter  http.Handler
	outboxRelay *outbox.Relay
}

func Build(conf *config.Config) *App {
	keycloakClient := keycloack.NewClient(
		conf.KeyCloak.URL,
		conf.KeyCloak.Realm,
		conf.KeyCloak.ClientID,
		conf.KeyCloak.ClientSecret,
		conf.KeyCloak.AdminUsername,
		conf.KeyCloak.AdminPassword,
	)

	schemaRegClient := schemaregistry.NewClient(conf.SchemaRegistry.BaseURL)
	schemaLocalReg, err := schemaregistry.NewLocalRegistry(context.Background(), schemaRegClient)
	if err != nil {
		panic(err)
	}

	schemaEncoder := schemaregistry.NewEncoder(schemaLocalReg)

	kClient, err := kafka.NewClient(conf.Kafka)
	if err != nil {
		panic(err)
	}

	if err := postgres.Migrate(conf.Postgres); err != nil {
		panic(err)
	}

	pgPool, err := postgres.NewPool(context.Background(), conf.Postgres)
	if err != nil {
		panic(err)
	}

	kProducer := kafka.NewProducer(kClient)
	outboxRepo := postgres.NewOutboxRepo(pgPool)
	outboxWriter := outbox.NewWriter(conf.Outbox, outboxRepo, schemaEncoder)
	outboxRelay := outbox.NewRelay(kProducer, outboxRepo, conf.Outbox, slog.Default())

	authService := services.NewAuth(keycloakClient, outboxWriter)

	deps := deliveryhttp.RouterDeps{
		AuthHandler: handlers.NewAuthHandler(
			authService,
			conf.KeyCloak.AccessTokenExpires,
			conf.KeyCloak.RefreshTokenExpires,
		),
	}
	httpRouter := deliveryhttp.NewRouter(&deps)

	return &App{Config: conf, httpRouter: httpRouter, outboxRelay: outboxRelay}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if app.outboxRelay != nil {
		go app.outboxRelay.Run(ctx)
	}

	err := http.ListenAndServe(app.Config.Addr, app.httpRouter)
	if err != nil {
		panic(err)
	}
}
