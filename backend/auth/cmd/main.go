package main

import (
	"log/slog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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
