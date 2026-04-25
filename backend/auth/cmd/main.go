package main

import (
	"log/slog"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/app"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
)

func main() {
	slog.Info("Service is starting")

	conf, err := config.Load()
	if err != nil {
		slog.Error("Error loading config", "error", err)
	}

	application := app.Build(conf)

	application.Run()

}
