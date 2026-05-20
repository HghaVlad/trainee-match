package main

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/elastic"
	"github.com/caarlos0/env/v11"
)

func main() {
	var cfg config.ElasticSearch

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	cl, err := elastic.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	err = elastic.Init(context.Background(), cl)
	if err != nil {
		panic(err)
	}
}
