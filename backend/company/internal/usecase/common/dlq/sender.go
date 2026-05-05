package dlq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

type producer interface {
	ProduceToDLQ(ctx context.Context, key, payload []byte, msg Message) error
}

type encoder interface {
	DLQToBytes(msg Message) ([]byte, error)
}

type Sender struct {
	producer producer
	encoder  encoder
	cfg      config.Kafka
}

func NewSender(cfg config.Kafka, producer producer, encoder encoder) *Sender {
	return &Sender{
		producer: producer,
		encoder:  encoder,
		cfg:      cfg,
	}
}

func (s *Sender) ToDLQ(
	ctx context.Context,
	eventID uuid.UUID,
	key, payload []byte,
	topic, eventType string,
	errMsg string,
) error {
	msg := Message{
		EventID:           eventID,
		Payload:           payload,
		OriginalEventType: eventType,
		OriginalTopic:     topic,
		LastError:         errMsg,
		FailedAt:          time.Now().UTC(),
	}

	b, err := s.encoder.DLQToBytes(msg)
	if err != nil {
		return fmt.Errorf("dlq encode: %w", err)
	}

	return s.producer.ProduceToDLQ(ctx, key, b, msg)
}
