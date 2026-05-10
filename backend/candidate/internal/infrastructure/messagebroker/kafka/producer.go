package kafka

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/outbox"

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

func (pr *Producer) ProduceOutBox(ctx context.Context, messages []outbox.Message) []outbox.ProduceResult {
	ctx, cancel := withDefaultTimeout(ctx, 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	results := make([]outbox.ProduceResult, len(messages))
	for i, m := range messages {
		results[i] = outbox.ProduceResult{MsgID: messages[i].ID}

		newRecord := kgo.Record{
			Key:     m.Key,
			Value:   m.Payload,
			Headers: mapHeadersWithEventType(m),
			Topic:   m.Topic,
		}

		wg.Add(1)

		pr.client.Produce(ctx, &newRecord, func(record *kgo.Record, err error) {
			defer wg.Done()

			if err == nil {
				now := time.Now()
				results[i].SentAt = &now
				return
			}

			results[i].Err = err

			pr.logger.Warn("outbox msg wasn't produced", "record", newRecord)

			var kErr *kerr.Error
			if errors.As(err, &kErr) && !kErr.Retriable {
				results[i].Unretryable = true
			}
		})
	}
	wg.Wait()
	return results
}

func withDefaultTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}

func mapHeadersWithEventType(msg outbox.Message) []kgo.RecordHeader {
	headers := make([]kgo.RecordHeader, 0, len(msg.Headers)+2)
	headers = append(headers, kgo.RecordHeader{
		Key: "event_type", Value: []byte(msg.EventType),
	})
	headers = append(headers, kgo.RecordHeader{
		Key: "event_id", Value: msg.ID[:],
	})

	for k, v := range msg.Headers {
		headers = append(headers, kgo.RecordHeader{
			Key: k, Value: []byte(v),
		})
	}

	return headers
}
