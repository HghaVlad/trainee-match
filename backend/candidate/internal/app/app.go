package app

import (
	"context"
	"errors"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/messagebroker/kafka"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/messagebroker/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/outbox"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"log/slog"
	"net/http"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	myhttp "github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate_by_user_id"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_skill"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
	Db     *pgxpool.Pool
	Relay  *outbox.Relay
}

func Build(conf *config.Config) (*App, error) {
	pgPool, err := postgres.Connect(context.Background(), &conf.Db)
	if err != nil {
		return nil, err
	}
	err = postgres.Migrate(conf.Db.GetPostgresURL())
	if err != nil {
		return nil, err
	}

	candidateRepo := repository.NewCandidateRepo(pgPool)
	resumeRepo := repository.NewResumeRepo(pgPool)
	skillRepo := repository.NewSkillRepo(pgPool)
	outboxRepository := repository.NewOutbox(pgPool, trmpgx.DefaultCtxGetter)

	schemaRegistryClient := schemaregistry.NewClient(conf.SchemaRegistry.BaseURL)
	schemalLocalRegistry, err := schemaregistry.NewLocalRegistry(context.Background(), schemaRegistryClient)
	if err != nil {
		return nil, err
	}
	encoder := schemaregistry.NewEncoder(schemalLocalRegistry)

	outboxWriter := outbox.NewWriter(conf.Outbox, outboxRepository, encoder)

	createCandidateUC := create_candidate.New(candidateRepo)
	updateCandidateUC := update_candidate.New(candidateRepo, outboxWriter)
	getCandidateByUserIdUC := get_candidate_by_user_id.New(candidateRepo)

	getResumeUC := get_resume.New(resumeRepo, candidateRepo)
	createResumeUC := create_resume.New(resumeRepo, skillRepo, candidateRepo)
	updateResumeUC := update_resume.New(resumeRepo, skillRepo, candidateRepo)

	getSkillUC := get_skill.New(skillRepo)

	candidateHandler := handlers.NewCandidate(createCandidateUC, updateCandidateUC, getCandidateByUserIdUC)
	resumeHandler := handlers.NewResume(createResumeUC, getResumeUC, updateResumeUC)
	skillHandler := handlers.NewSkill(getSkillUC)
	authMiddleware := auth.NewMiddleware(conf.JWKUrl)

	router := myhttp.NewRouter(myhttp.NewRouterDeps(authMiddleware, candidateHandler, resumeHandler, skillHandler))

	httpServer := &http.Server{
		Addr:         conf.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	kafkaClient, err := kafka.NewClient(conf.Kafka)
	if err != nil {
		return nil, err
	}

	kafkaLogger := slog.Logger{}
	kafkaProducer := kafka.NewProducer(kafkaClient, conf.Kafka, &kafkaLogger)
	trManager := manager.Must(trmpgx.NewFactory(pgPool))
	relayLogger := slog.Logger{}
	outboxRelay := outbox.NewRelay(outboxRepository, kafkaProducer, conf.Outbox, &relayLogger, trManager)

	return &App{
		server: httpServer,
		Db:     pgPool,
		Relay:  outboxRelay,
	}, nil
}

func (app *App) Run() error {

	slog.Info("Server started")
	err := app.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http listening server err", "error", err)
	}

	go app.Relay.Run(context.Background())

	return err
}

func (app *App) Shutdown(ctx context.Context) {
	err := app.server.Shutdown(ctx)
	if err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("Server stopped")
	app.Db.Close()
}
