package kafka

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
)

type Producer struct {
	client *kgo.Client
	config config.Kafka
	logger *slog.Logger
}

func NewProducer(client *kgo.Client, config config.Kafka, logger *slog.Logger) *Producer {
	return &Producer{client: client, config: config, logger: logger}
}

func (pr *Producer) Close() {
	if pr.client != nil {
		pr.client.Close()
	}
}

func (pr *Producer) ProduceOutBox(ctx context.Context, messages []OutBox) ([]OutBox, error) {}
