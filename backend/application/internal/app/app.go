package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/resumearchived"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/messaging/kafka"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/messaging/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/utils/hash"
	apphttp "github.com/HghaVlad/trainee-match/backend/application/internal/transport/http"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/dynamics"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/summary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/candidatehistory"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/getcandidateview"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/gethrview"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/hrhistory"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/hrupdstatus"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listhrsummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/withdraw"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/dlq"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/eventhandler"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/candidateupserted"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/companydeleted"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/companymemberadded"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/companymemberremoved"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/companymodupd"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/companyupdated"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/resumedeleted"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/resumeupserted"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/vacancyarchived"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/vacancymodupd"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/vacancypublished"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/vacancyupdated"
)

type App struct {
	httpServer    *http.Server
	pgDB          *pgxpool.Pool
	logger        *slog.Logger
	kafkaConsumer *kafka.Consumer
	kafkaProducer *kafka.Producer
}

//nolint:funlen // app wiring
func Build(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	pgDB, err := postgres.ConnectPgxPoolWithLogger(ctx, cfg.DB, logger)
	if err != nil {
		return nil, err
	}

	txManager := manager.Must(trmpgx.NewDefaultFactory(pgDB))
	txGetter := trmpgx.DefaultCtxGetter

	appRepo := repository.NewApplication(pgDB, txGetter)
	appSnapRepo := repository.NewAppSnapshot(pgDB, txGetter)
	appStatusHistoryRepo := repository.NewAppStatusHistoryRepo(pgDB, txGetter)
	resumeProjRepo := repository.NewResumeProjection(pgDB, txGetter)
	vacProjRepo := repository.NewVacancyProjection(pgDB, txGetter)
	candProjRepo := repository.NewCandidateProjection(pgDB, txGetter)
	compMemProjRepo := repository.NewCompanyMemberProjection(pgDB, txGetter)

	appSnapHasher := hash.NewAppSnapshotHasher()

	schemaRegistryClient := schemaregistry.NewClient(cfg.SchemaRegistry)
	localRegistry, err := schemaregistry.NewLocalRegistry(ctx, schemaRegistryClient)
	if err != nil {
		return nil, err
	}
	decoder := schemaregistry.NewDecoder(localRegistry)
	encoder := schemaregistry.NewEncoder(localRegistry)

	applyUC := apply.NewUsecase(
		appRepo,
		resumeProjRepo,
		candProjRepo,
		vacProjRepo,
		appSnapRepo,
		appStatusHistoryRepo,
		appSnapHasher,
		txManager,
	)
	listCandidateAppsUC := listcandidatesummary.NewUsecase(appRepo)
	getCandidateView := getcandidateview.NewUsecase(appRepo)
	getCandiAppHistory := candidatehistory.NewUsecase(appStatusHistoryRepo)
	withdrawApp := withdraw.NewUsecase(appRepo, appStatusHistoryRepo, txManager)

	listHrApps := listhrsummary.NewUsecase(appRepo, compMemProjRepo, vacProjRepo)
	getHrDetailedView := gethrview.NewUsecase(appRepo)
	hrUpdateStatus := hrupdstatus.NewUsecase(appRepo, appStatusHistoryRepo, txManager)
	getHistoryHrView := hrhistory.NewUsecase(appStatusHistoryRepo)

	analyticsSummary := summary.NewUsecase(appRepo, compMemProjRepo, vacProjRepo)
	dynamicsDashboard := dynamics.NewUsecase(appStatusHistoryRepo, compMemProjRepo, vacProjRepo)

	authMiddleware, err := middleware.NewAuthMiddleware(ctx, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	deps := &handlers.Deps{
		Apply:               applyUC,
		ListCandidateApps:   listCandidateAppsUC,
		GetCandidateViewUC:  getCandidateView,
		GetCandiStatHistory: getCandiAppHistory,
		Withdraw:            withdrawApp,
		ListHrApps:          listHrApps,
		GetHrDetailedView:   getHrDetailedView,
		HrUpdateStatus:      hrUpdateStatus,
		GetHistoryHrView:    getHistoryHrView,
		AnalyticsSummary:    analyticsSummary,
		DynamicsDashboard:   dynamicsDashboard,
		Logger:              logger,
	}

	hand := handlers.NewHandler(deps)

	// Event usecases
	resumeUpserted := resumeupserted.NewUsecase(resumeProjRepo)
	resumeDeleted := resumedeleted.NewUsecase(resumeProjRepo)
	resumeArchived := resumearchived.NewUsecase(resumeProjRepo)
	candidateUpserted := candidateupserted.NewUsecase(candProjRepo)
	companyUpdated := companyupdated.NewUsecase(vacProjRepo)
	companyDeleted := companydeleted.NewUsecase(compMemProjRepo, vacProjRepo, appRepo, appStatusHistoryRepo, txManager)
	companyMemberAdded := companymemberadded.NewUsecase(compMemProjRepo)
	companyMemberRemoved := companymemberremoved.NewUsecase(compMemProjRepo)
	vacancyPublished := vacancypublished.NewUsecase(vacProjRepo)
	vacancyArchived := vacancyarchived.NewUsecase(vacProjRepo)
	vacancyUpdated := vacancyupdated.NewUsecase(vacProjRepo)
	vacancyModUpd := vacancymodupd.NewUsecase(vacProjRepo)
	companyModUpd := companymodupd.NewUsecase(vacProjRepo)

	// Kafka Producer
	kafkaProducerClient, err := kafka.NewClientForProducer(cfg.Kafka)
	if err != nil {
		return nil, err
	}
	kafkaProducer := kafka.NewProducer(kafkaProducerClient, cfg.Kafka)
	dlqSender := dlq.NewSender(kafkaProducer, encoder)

	// Kafka Consumer
	eventHandler := eventhandler.NewHandler(
		decoder,
		cfg.KafkaHandling,
		logger,
		dlqSender,
		resumeUpserted,
		resumeDeleted,
		resumeArchived,
		candidateUpserted,
		companyUpdated,
		companyDeleted,
		companyMemberAdded,
		companyMemberRemoved,
		vacancyPublished,
		vacancyArchived,
		vacancyUpdated,
		vacancyModUpd,
		companyModUpd,
	)

	consumer := kafka.NewConsumer(eventHandler, logger)
	kafkaConsumerClient, err := kafka.NewClientForConsumer(cfg.Kafka, consumer)
	if err != nil {
		return nil, err
	}
	consumer.Client = kafkaConsumerClient

	router := apphttp.NewRouter(hand, authMiddleware, logger)

	httpServer := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		httpServer:    httpServer,
		logger:        logger,
		pgDB:          pgDB,
		kafkaConsumer: consumer,
		kafkaProducer: kafkaProducer,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		app.logger.InfoContext(ctx, "starting http server", "addr", app.httpServer.Addr)

		if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		app.logger.Info("starting the kafka consumer")
		app.kafkaConsumer.Poll(ctx)
		app.logger.Info("finishing the kafka consumer")
		return nil
	})

	return g.Wait()
}

func (app *App) Shutdown(ctx context.Context) {
	if err := app.httpServer.Shutdown(ctx); err != nil {
		app.logger.InfoContext(ctx, "http server graceful shutdown fail", "error", err)

		if cerr := app.httpServer.Close(); cerr != nil {
			app.logger.InfoContext(ctx, "http server close fail", "error", cerr)
		}
	}

	app.pgDB.Close()

	app.kafkaProducer.Shutdown()
	app.kafkaConsumer.Shutdown()

	app.logger.InfoContext(ctx, "app gracefully stopped")
}
