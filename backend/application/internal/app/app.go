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

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres/repository"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/utils/hash"
	apphttp "github.com/HghaVlad/trainee-match/backend/application/internal/transport/http"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/handlers"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/getcandidateview"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
)

type App struct {
	httpServer *http.Server
	pgDB       *pgxpool.Pool
	logger     *slog.Logger
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

	appSnapHasher := hash.NewAppSnapshotHasher()

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

	authMiddleware, err := middleware.NewAuthMiddleware(ctx, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	deps := &handlers.Deps{
		Apply:              applyUC,
		ListCandidateApps:  listCandidateAppsUC,
		GetCandidateViewUC: getCandidateView,
		Logger:             logger,
	}

	hand := handlers.NewHandler(deps)

	router := apphttp.NewRouter(hand, authMiddleware, logger)

	httpServer := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		httpServer: httpServer,
		logger:     logger,
		pgDB:       pgDB,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	app.logger.InfoContext(ctx, "starting http server", "addr", app.httpServer.Addr)

	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.logger.InfoContext(ctx, "http server stopped")
	return nil
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
