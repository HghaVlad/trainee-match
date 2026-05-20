package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/dlq"
)

type Producer struct {
	client *kgo.Client
	cfg    config.Kafka
}

func NewProducer(client *kgo.Client, cfg config.Kafka) *Producer {
	return &Producer{client: client, cfg: cfg}
}

func (pr *Producer) ProduceMessages(ctx context.Context, message []byte) error {
	_ = ctx
	_ = message
	return fmt.Errorf("ProduceMessages is not implemented")
}

func (pr *Producer) ProduceDLQ(ctx context.Context, message dlq.Message, key, value []byte) error {
	record := kgo.Record{
		Topic: pr.cfg.DLQTopic,
		Key:   key,
		Value: value,

		Headers: []kgo.RecordHeader{
			{Key: "event_type", Value: []byte("dlq")},
			{Key: "event_id", Value: []byte(message.EventID.String())},
			{Key: "original_event_type", Value: []byte(message.OriginalEventType)},
			{Key: "original_topic", Value: []byte(message.OriginalTopic)},
			{Key: "error", Value: []byte(message.LastError)},
			{Key: "failed_at", Value: []byte(message.FailedAt.String())},
		},
	}

	err := pr.client.ProduceSync(ctx, &record).FirstErr()
	if err != nil {
		return fmt.Errorf("failed to produce DLQ message: %w", err)
	}

	return nil
}

func (pr *Producer) Shutdown() {
	pr.client.Close()
}
