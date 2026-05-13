package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserCreatedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	UserID     uuid.UUID `avro:"user_id"`
	Username   string    `avro:"username"`
	Email      string    `avro:"email"`
	Role       string    `avro:"role"`
	OccurredAt time.Time `avro:"occurred_at"`
}
