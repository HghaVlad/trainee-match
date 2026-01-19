package main

import (
	"log"

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
		log.Fatal(err)
	}

	myApp := app.Build(conf)
	myApp.Run()

	log.Println("Service has started")
}
