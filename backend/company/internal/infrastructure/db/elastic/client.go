package elastic

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

func NewClient(cfg config.ElasticSearch) (*elasticsearch.Client, error) {
	cl, err := elasticsearch.New(
		elasticsearch.WithAddresses(cfg.Addrs...),
		elasticsearch.WithBasicAuth(cfg.User, cfg.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("elastic search create client: %v", err)
	}

	return cl, nil
}
