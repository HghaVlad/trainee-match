package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type HrAppSummary struct {
	AppID        uuid.UUID
	Status       application.Status
	VacancyID    uuid.UUID
	VacancyTitle string
	AppSnap      HrAppSnapSummary
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type HrAppSnapSummary struct {
	Email     string
	FullName  string
	Telegram  *string
	CreatedAt time.Time
}
