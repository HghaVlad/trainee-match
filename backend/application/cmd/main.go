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
	os.Exit(run())
}

func run() int {
	logger := slog.Default()
	logger.Info("starting application service...")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("error loading config", "error", err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	myApp, err := app.Build(ctx, *cfg, logger)
	if err != nil {
		logger.Error("error building app", "error", err)
		return 1
	}

	runErrCh := make(chan error, 1)

	go func() {
		err = myApp.Run(ctx)
		if err != nil {
			logger.Error("app stopped with error", "error", err)
			runErrCh <- err
		}
	}()

	exitCode := 0

	select {
	case <-ctx.Done():
	case <-runErrCh:
		exitCode = 1
		stop()
	}

	logger.Info("gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myApp.Shutdown(shutdownCtx)
	logger.Info("app stopped")
	return exitCode
}
