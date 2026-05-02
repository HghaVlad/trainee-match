package kafka

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
)

type Producer struct {
	client *kgo.Client
	logger *slog.Logger
}

func NewProducer(client *kgo.Client, logger *slog.Logger) *Producer {
	return &Producer{
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
