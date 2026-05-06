package company

import (
	"time"

	"github.com/google/uuid"
)

type UpdatedEvent struct {
	EventID     uuid.UUID `avro:"event_id"`
	CompanyID   uuid.UUID `avro:"company_id"`
	CompanyName string    `avro:"company_name"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

type DeletedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	CompanyID  uuid.UUID `avro:"company_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}
