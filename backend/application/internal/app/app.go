package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	http2 "github.com/HghaVlad/trainee-match/backend/application/internal/delivery/http"
	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/db/postgres"
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

	deps := http2.NewRouterDeps()
	router := http2.NewRouter(deps)

	httpServer := &http.Server{
		Addr:         cfg.Http.Addr,
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
