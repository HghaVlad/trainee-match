package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/app"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
)

// @title Trainee Match: Candidate Service API
// @version 1.0
// @description API for managing candidate profiles
// @host 0.0.0.0:8081
// @BasePath /api/v1
// @schemes http https
func main() {
	logger := slog.Default()
	logger.Info("Service is starting")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := config.Load()
	if err != nil {
		logger.Error("Error loading config", "error", err)
	}
	logger.Debug("Config loaded", "config", conf)

	myApp, err := app.Build(conf, logger)
	if err != nil {
		logger.Error("Error building app", "error", err)
		return
	}
	logger.Info("App built")

	errChan := make(chan error)
	go func() {
		errChan <- myApp.Run()
	}()
	logger.Info("App is run")

	select {
	case <-ctx.Done():
	case <-errChan:
	}
	logger.Info("Shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myApp.Shutdown(shutdownCtx)
	logger.Info("Service stopped")
}
