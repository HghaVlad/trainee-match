package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/eventhandler"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/projection/resumeupserted"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/messaging/schemaregistry"

	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/messaging/kafka"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres/repository"
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
)

type App struct {
	httpServer    *http.Server
	pgDB          *pgxpool.Pool
	logger        *slog.Logger
	kafkaConsumer *kafka.Consumer
}

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

	//Event usecases
	resumeUpserted := resumeupserted.NewUsecase(resumeProjRepo)

	// Kafka
	eventHandler := eventhandler.NewHandler(decoder, resumeUpserted)

	consumer := kafka.NewConsumer(eventHandler)
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

	if err := g.Wait(); err != nil {
		return err
	}

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

	app.logger.InfoContext(ctx, "app gracefully stopped")
}
