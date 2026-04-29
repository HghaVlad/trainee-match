package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID           uuid.UUID
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

	// in case smth goes wrong
	LastError     *string
	NextAttemptAt time.Time
	FailedAt      *time.Time
}

type Status string

const (
	StatusPending Status = "pending"
	StatusSent    Status = "sent"
	StatusFailed  Status = "dead"
)
