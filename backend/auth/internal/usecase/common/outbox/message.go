package outbox

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a single outbox event to be published to Kafka.
type Message struct {
	ID           uuid.UUID
	AggregateID  uuid.UUID
	AggregateSeq int64
	Topic        string
	Key          []byte
	Payload      []byte
	Headers      map[string]string
	SchemaID     int
	EventType    EventType
	Status       Status
	AttemptCount int
	MaxAttempts  int
	CreatedAt    time.Time
	SentAt       *time.Time

	LastError     *string
	NextAttemptAt time.Time
	FailedAt      *time.Time
}

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
)
