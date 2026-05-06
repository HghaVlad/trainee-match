package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/app"
	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
)

// @title Trainee Match: Company Service API
// @version 1.0
// @description company microservice
// @BasePath /api/v1
// @schemes http https
func main() {
	os.Exit(run())
}

// actual logic of main,
// returns exit code, and all defers work normally
func run() int {
	lgr := logger.NewSlogLogger()
	lgr.Info("Service is starting...")

	cfg, err := config.Load()
	if err != nil {
		lgr.Error("couldn't load config", "error", err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	myApp, err := app.Build(ctx, cfg, lgr)
	if err != nil {
		lgr.Error("couldn't build app", "error", err)
		return 1
	}

	runErrCh := make(chan error, 1)

	go func() {
		err = myApp.Run(ctx)
		if err != nil {
			lgr.Error("app stopped with error", "error", err)
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

	lgr.Info("Gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myApp.Shutdown(shutdownCtx)
	lgr.Info("server stopped")
	return exitCode
}
