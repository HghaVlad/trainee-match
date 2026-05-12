package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type HrDetailedView struct {
	AppID          uuid.UUID
	Status         application.Status
	VacancyID      uuid.UUID
	VacancyTitle   string
	CoverLetter    *string
	Snapshot       ApplicationSnapshot
	StatusHistory  []StatusChangeHrFullView
	CreatedAt      time.Time
	UpdatedAt      time.Time
	AllowedActions []AllowedAction
}

type StatusChangeHrFullView struct {
	Status          application.Status `json:"status"`
	CreatedAt       time.Time          `json:"created_at"`
	Comment         *string            `json:"comment"`
	ChangedByRole   application.Actor  `json:"changed_by_role"`
	ChangedByUserID *uuid.UUID         `json:"changed_by_user_id"`
}
