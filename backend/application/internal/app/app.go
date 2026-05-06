package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	http2 "github.com/HghaVlad/trainee-match/backend/application/internal/delivery/http"
)

type App struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func Build(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
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
	}, nil
}

func (app *App) Run(_ context.Context) error {
	app.logger.Info("starting http server", "addr", app.httpServer.Addr)
	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.logger.Info("http server stopped")
	return nil
}

func (app *App) Shutdown(ctx context.Context) {
	if err := app.httpServer.Shutdown(ctx); err != nil {
		app.logger.Info("http server graceful shutdown fail", "error", err)

		if cerr := app.httpServer.Close(); cerr != nil {
			app.logger.Info("server close fail", cerr)
		}

		return
	}

	app.logger.Info("app gracefully stopped")
}
