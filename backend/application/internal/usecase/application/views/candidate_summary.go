package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type CandidateAppSummary struct {
	AppID        uuid.UUID
	Status       application.Status
	VacancyID    uuid.UUID
	VacancyTitle string
	CompanyID    uuid.UUID
	CompanyName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
