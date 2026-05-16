package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CompanyMember struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Role      string
}

var (
	ErrCompanyMemberNotFound = errors.New("company member not found")
	ErrCompanyNotFound       = errors.New("company not found")
)

type CompanyMemberAddedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	UserID     uuid.UUID `avro:"user_id"`
	CompanyID  uuid.UUID `avro:"company_id"`
	Role       string    `avro:"role"`
	OccurredAt time.Time `avro:"occurred_at"`
}

func (ev CompanyMemberAddedEvent) ToCompanyMember() CompanyMember {
	return CompanyMember{
		UserID:    ev.UserID,
		CompanyID: ev.CompanyID,
		Role:      ev.Role,
	}
}

type CompanyMemberRemovedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	UserID     uuid.UUID `avro:"user_id"`
	CompanyID  uuid.UUID `avro:"company_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}

type CompanyUpdatedEvent struct {
	EventID     uuid.UUID `avro:"event_id"`
	CompanyID   uuid.UUID `avro:"company_id"`
	CompanyName string    `avro:"company_name"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

type CompanyDeletedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	CompanyID  uuid.UUID `avro:"company_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}
