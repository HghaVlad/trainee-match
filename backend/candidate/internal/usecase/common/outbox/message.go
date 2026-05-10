package outbox

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID           uuid.UUID
	AggregateID  uuid.UUID
	AggregateSeq int64
	Topic        string
	EventType    EventType
	SchemaID     int
	Headers      map[string]string
	Key          []byte
	Payload      []byte
	Status       Status
	AttemptCount int64
	CreatedAt    time.Time
	SentAt       *time.Time

	// in case smth goes wrong
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
