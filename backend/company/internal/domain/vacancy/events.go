package vacancy

import (
	"time"

	"github.com/google/uuid"
)

type PublishedEvent struct {
	EventID     uuid.UUID `avro:"event_id"`
	VacancyID   uuid.UUID `avro:"vacancy_id"`
	Title       string    `avro:"title"`
	CompanyID   uuid.UUID `avro:"company_id"`
	CompanyName string    `avro:"company_name"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

type UpdatedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	VacancyID  uuid.UUID `avro:"vacancy_id"`
	Title      string    `avro:"title"`
	OccurredAt time.Time `avro:"occurred_at"`
}

type ArchivedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	VacancyID  uuid.UUID `avro:"vacancy_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}
