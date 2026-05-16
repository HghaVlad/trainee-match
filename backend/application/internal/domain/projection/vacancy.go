package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type VacancyStatus string

const (
	VacancyStatusPublished VacancyStatus = "published"
	VacancyStatusArchived  VacancyStatus = "archived"
)

type Vacancy struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	CompanyName string
	Title       string
	Status      VacancyStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type VacancyPublishedEvent struct {
	EventID     uuid.UUID `avro:"event_id"`
	VacancyID   uuid.UUID `avro:"vacancy_id"`
	Title       string    `avro:"title"`
	CompanyID   uuid.UUID `avro:"company_id"`
	CompanyName string    `avro:"company_name"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

func (ev VacancyPublishedEvent) ToVacancy() Vacancy {
	occurredAt := ev.OccurredAt
	return Vacancy{
		ID:          ev.VacancyID,
		CompanyID:   ev.CompanyID,
		CompanyName: ev.CompanyName,
		Title:       ev.Title,
		Status:      VacancyStatusPublished,
		CreatedAt:   &occurredAt,
		UpdatedAt:   &occurredAt,
	}
}

type VacancyUpdatedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	VacancyID  uuid.UUID `avro:"vacancy_id"`
	Title      string    `avro:"title"`
	OccurredAt time.Time `avro:"occurred_at"`
}

type VacancyArchivedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	VacancyID  uuid.UUID `avro:"vacancy_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}

var (
	ErrVacancyNotFound = errors.New("vacancy not found")
)
