package domain

import (
	"github.com/google/uuid"
	"time"
)

type ResumeUpdatedEvent struct {
	EventID    uuid.UUID `json:"event_id"`
	Resume     Resume
	OccurredAt time.Time `json:"occurred_at"`
}
