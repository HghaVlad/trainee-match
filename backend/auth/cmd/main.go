package main

import (
	"github.com/HghaVlad/trainee-match/backend/auth/internal/app"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
	"log/slog"
)

func main() {
	slog.Info("Service is starting")

	conf, err := config.Load()
	if err != nil {
		slog.Error("Error loading config", err)
	}

	application := app.Build(conf)

	application.Run()

}
