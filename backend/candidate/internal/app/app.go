package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/messagebroker/kafka"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/infrastructure/messagebroker/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/common/outbox"

	"github.com/jackc/pgx/v5/pgxpool"

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
)

type App struct {
	server      *http.Server
	Db          *pgxpool.Pool
	Relay       *outbox.Relay
	relayCancel context.CancelFunc
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

	trManager := manager.Must(trmpgx.NewFactory(pgPool))

	schemaRegistryClient := schemaregistry.NewClient(conf.SchemaRegistry.BaseURL)
	schemalLocalRegistry, err := schemaregistry.NewLocalRegistry(context.Background(), schemaRegistryClient)
	if err != nil {
		return nil, err
	}
	encoder := schemaregistry.NewEncoder(schemalLocalRegistry)

	outboxWriter := outbox.NewWriter(conf.Outbox, outboxRepository, encoder)

	createCandidateUC := create_candidate.New(candidateRepo, outboxWriter, trManager)
	updateCandidateUC := update_candidate.New(candidateRepo, outboxWriter, trManager)
	getCandidateByUserIdUC := get_candidate_by_user_id.New(candidateRepo)

	getResumeUC := get_resume.New(resumeRepo, candidateRepo)
	createResumeUC := create_resume.New(resumeRepo, skillRepo, candidateRepo, outboxWriter, trManager)
	updateResumeUC := update_resume.New(resumeRepo, skillRepo, candidateRepo, outboxWriter, trManager)

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

	kafkaLogger := slog.New(slog.NewTextHandler(log.Writer(), nil))
	kafkaProducer := kafka.NewProducer(kafkaClient, conf.Kafka, kafkaLogger)

	relayLogger := slog.New(slog.NewTextHandler(log.Writer(), nil))
	outboxRelay := outbox.NewRelay(outboxRepository, kafkaProducer, conf.Outbox, relayLogger, trManager)

	return &App{
		server: httpServer,
		Db:     pgPool,
		Relay:  outboxRelay,
	}, nil
}

func (app *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	app.relayCancel = cancel
	go app.Relay.Run(ctx)

	slog.Info("Server started")
	err := app.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http listening server err", "error", err)
	}

	return err
}

func (app *App) Shutdown(ctx context.Context) {
	if app.relayCancel != nil {
		app.relayCancel()
	}
	err := app.server.Shutdown(ctx)
	if err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("Server stopped")
	app.Db.Close()
}
