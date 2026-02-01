package main

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/app"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

// @title Trainee Match: Candidate Service API
// @version 1.0
// @description API for managing candidate profiles
// @host 0.0.0.0:8081
// @BasePath /api/v1
// @schemes http https
func main() {
	slog.Info("Service is starting")
	slog.SetLogLoggerLevel(-100)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := config.Load()
	if err != nil {
		slog.Error("Error loading config", err)
	}
	slog.Debug("Config loaded", "config", conf)

	myApp, err := app.Build(conf)
	if err != nil {
		slog.Error("Error building app", err)
		return
	}
	slog.Info("App built")

	errChan := make(chan error)
	go func() {
		errChan <- myApp.Run()
	}()

	select {
	case <-ctx.Done():
	case <-errChan:
	}
	slog.Info("Shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myApp.Shutdown(shutdownCtx)
	slog.Info("Service stopped")
}
