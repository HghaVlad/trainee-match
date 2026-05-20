package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/admin/newadmin"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
	deliveryhttp "github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/keycloack"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/message_broker/kafka"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/infra/message_broker/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/common/outbox"
	getuser "github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/get_user_me"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/login"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/logout"
	refreshtoken "github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/refresh_token"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/register"
)

type App struct {
	Config      *config.Config
	httpServer  *http.Server
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
	trManager := manager.Must(trmpgx.NewFactory(pgPool))
	kProducer := kafka.NewProducer(kClient)
	outboxRepo := postgres.NewOutboxRepo(pgPool, trmpgx.DefaultCtxGetter)
	outboxWriter := outbox.NewWriter(conf.Outbox, outboxRepo, schemaEncoder)
	outboxRelay := outbox.NewRelay(kProducer, outboxRepo, conf.Outbox, slog.Default(), trManager)

	authRegisterUC := register.New(keycloakClient, outboxWriter)
	authLoginUC := login.New(keycloakClient)
	authLogoutUC := logout.New(keycloakClient)
	authRefreshUC := refreshtoken.New(keycloakClient)
	authGetMeUC := getuser.New(keycloakClient)

	adminNew := newadmin.NewUseCase(keycloakClient)

	deps := deliveryhttp.RouterDeps{
		AuthHandler: handlers.NewAuthHandler(
			authRegisterUC,
			authLoginUC,
			authLogoutUC,
			authRefreshUC,
			authGetMeUC,
			conf.KeyCloak.AccessTokenExpires,
			conf.KeyCloak.RefreshTokenExpires,
		),
		AdminHandler: handlers.NewAdmin(
			adminNew,
		),
	}
	handler := deliveryhttp.NewRouter(&deps)

	server := &http.Server{
		Addr:         conf.Addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{Config: conf, httpServer: server, outboxRelay: outboxRelay}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if app.outboxRelay != nil {
		go app.outboxRelay.Run(ctx)
	}

	err := app.httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
