package main

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/elastic"
)

func main() {
	cfg := config.ElasticSearch{
		Addrs: []string{"http://localhost:9200"},
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
