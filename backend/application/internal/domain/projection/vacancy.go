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

type ModerationStatus string

const (
	ModerationStatusOK     ModerationStatus = "ok"
	ModerationStatusHidden ModerationStatus = "hidden"
)

type Vacancy struct {
	ID            uuid.UUID
	CompanyID     uuid.UUID
	CompanyName   string
	Title         string
	Status        VacancyStatus
	ModStatus     ModerationStatus
	CompModStatus ModerationStatus
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func (v *Vacancy) IsApplyable() error {
	if v.Status != VacancyStatusPublished {
		return ErrVacancyNotPublished
	}

	if v.ModStatus != ModerationStatusOK || v.CompModStatus != ModerationStatusOK {
		return ErrVacancyBadModStatus
	}

	return nil
}

var (
	ErrVacancyNotFound         = errors.New("vacancy not found")
	ErrVacancyNotPublished = errors.New("vacancy must be published")
	ErrVacancyBadModStatus = errors.New("vacancy invalid moderation status")
)

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
		ID:            ev.VacancyID,
		CompanyID:     ev.CompanyID,
		CompanyName:   ev.CompanyName,
		Title:         ev.Title,
		Status:        VacancyStatusPublished,
		CreatedAt:     &occurredAt,
		UpdatedAt:     &occurredAt,
		ModStatus:     ModerationStatusOK,
		CompModStatus: ModerationStatusOK,
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

type VacancyModerationUpdatedEvent struct {
	EventID    uuid.UUID        `avro:"event_id"`
	VacancyID  uuid.UUID        `avro:"vacancy_id"`
	ModStatus  ModerationStatus `avro:"moderation_status"`
	OccurredAt time.Time        `avro:"occurred_at"`
}

type CompanyModerationUpdatedEvent struct {
	EventID    uuid.UUID        `avro:"event_id"`
	CompanyID  uuid.UUID        `avro:"company_id"`
	ModStatus  ModerationStatus `avro:"moderation_status"`
	OccurredAt time.Time        `avro:"occurred_at"`
}
