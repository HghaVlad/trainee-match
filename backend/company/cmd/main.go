package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/app"
	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

// @title Trainee Match: Company Service API
// @version 1.0
// @description company microservice
// @BasePath /api/v1
// @schemes http https
func main() {
	log.Println("Service is starting...")

	conf, err := config.Load()
	if err != nil {
		log.Fatal("config load err: ", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	myApp, err := app.Build(ctx, conf)
	if err != nil {
		stop()
		//nolint:gocritic // stop() is called right above
		log.Fatal("app build err: ", err)
	}

	go myApp.Run()

	log.Println("http listening on ", conf.HTTP.Addr)

	<-ctx.Done()
	log.Println("Gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myApp.Shutdown(shutdownCtx)
	log.Println("server stopped")
}
