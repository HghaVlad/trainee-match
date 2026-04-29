package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	http2 "github.com/HghaVlad/trainee-match/backend/application/internal/delivery/http"
)

type App struct {
	Server *http.Server
}

func Build(conf *config.Config) (*App, error) {
	deps := http2.NewRouterDeps()
	router := http2.NewRouter(deps)

	httpServer := &http.Server{
		Addr:         conf.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{Server: httpServer}, nil
}

func (app *App) Run() error {
	if app == nil || app.Server == nil {
		return fmt.Errorf("no server configured, nothing to run")
	}

	slog.Info("starting server with", "server address", app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	slog.Info("server stopped")
	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	if app == nil || app.Server == nil {
		return nil
	}

	if err := app.Server.Shutdown(ctx); err != nil {
		slog.Debug("graceful shutdown failed: %v", err)
		if cerr := app.Server.Close(); cerr != nil {
			slog.Debug("server close failed: %v", cerr)
		}
		return err
	}

	slog.Info("server gracefully stopped")
	return nil
}
