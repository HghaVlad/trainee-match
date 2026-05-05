package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/dlq"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
)

type Producer struct {
	client *kgo.Client
	cfg    config.Kafka
	logger *slog.Logger
}

func NewProducer(cfg config.Kafka, client *kgo.Client, logger *slog.Logger) *Producer {
	return &Producer{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

func (c *Producer) Close() {
	if c.client == nil {
		return
	}

	c.client.Close()
}

func (c *Producer) ProduceOutbox(ctx context.Context, msgs []outbox.Message) []outbox.ProduceResult {
	wg := new(sync.WaitGroup)

	results := make([]outbox.ProduceResult, len(msgs))

	for i := range msgs {
		results[i].MsgID = msgs[i].ID

		headers := mapHeadersWithEventType(msgs[i])

		record := &kgo.Record{
			Topic:   msgs[i].Topic,
			Key:     msgs[i].Key,
			Value:   msgs[i].Payload,
			Headers: headers,
		}

		wg.Add(1)

		c.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
			defer wg.Done()

			if err == nil {
				now := time.Now().UTC()
				results[i].SentAt = &now
				return
			}

			results[i].Err = err

			c.logger.Warn("outbox msg wasn't produced",
				"topic", record.Topic, "err", err)

			var kErr *kerr.Error
			if errors.As(err, &kErr) && !kErr.Retriable {
				results[i].Unretryable = true
			}
		})
	}

	wg.Wait()
	return results
}

func (c *Producer) ProduceToDLQ(ctx context.Context, key, payload []byte, msg dlq.Message) error {
	record := &kgo.Record{
		Topic: c.cfg.DLQTopic,
		Key:   key,
		Value: payload,
	}

	// duplicate to headers for better observability
	headers := []kgo.RecordHeader{
		{Key: "event_type", Value: []byte("dlq")},
		{Key: "original_event_type", Value: []byte(msg.OriginalEventType)},
		{Key: "original_topic", Value: []byte(msg.OriginalTopic)},
		{Key: "error", Value: []byte(msg.LastError)},
		{Key: "failed_at", Value: []byte(msg.FailedAt.String())},
	}
	record.Headers = headers

	err := c.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("dlq produce: %w", err)
	}

	return nil
}

func mapHeadersWithEventType(msg outbox.Message) []kgo.RecordHeader {
	result := make([]kgo.RecordHeader, 0, len(msg.Headers)+2)

	result = append(result, kgo.RecordHeader{
		Key:   "event_type",
		Value: []byte(msg.EventType),
	}, kgo.RecordHeader{
		Key:   "event_id",
		Value: msg.ID[:],
	})

	for key, value := range msg.Headers {
		result = append(result, kgo.RecordHeader{
			Key:   key,
			Value: []byte(value),
		})
	}

	return result
}
