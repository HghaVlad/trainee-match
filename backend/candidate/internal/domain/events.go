package domain

import "github.com/google/uuid"

type ResumeUpdatedEvent struct {
	EventID uuid.UUID `json:"event_id"`
	Resume  Resume
}
