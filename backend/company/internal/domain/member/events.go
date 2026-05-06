package member

import (
	"time"

	"github.com/google/uuid"
)

type AddedEvent struct {
	EventID    uuid.UUID   `avro:"event_id"`
	UserID     uuid.UUID   `avro:"user_id"`
	CompanyID  uuid.UUID   `avro:"company_id"`
	Role       CompanyRole `avro:"role"`
	OccurredAt time.Time   `avro:"occurred_at"`
}

type RemovedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	UserID     uuid.UUID `avro:"user_id"`
	CompanyID  uuid.UUID `avro:"company_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}
