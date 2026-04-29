package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/app"
	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
)

func main() {
	cfg := config.Load()
	slog.Info("starting application", "config", cfg)
	myApp, err := app.Build(cfg)
	if err != nil {
		slog.Error("error building app", "err", err)
		return
	}

	runErr := make(chan error, 1)
	go func() {
		runErr <- myApp.Run()
	}()

	// wait for signal or run error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		slog.Info("received signal, initiating shutdown", "signal", sig)
	case err = <-runErr:
		if err != nil {
			slog.Error("server run error, initiating shutdown", "err", err)
		} else {
			slog.Info("server stopped")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = myApp.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown gracefully", "err", err)
	}
}
