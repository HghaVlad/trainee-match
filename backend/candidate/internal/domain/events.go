package domain

import (
	"time"

	"github.com/google/uuid"
)

type ResumeUpdatedEvent struct {
	EventID    uuid.UUID `json:"event_id"`
	Resume     Resume
	OccurredAt time.Time `json:"occurred_at"`
}
