package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type AllowedAction string

const (
	AllowedActionWithdraw AllowedAction = "withdraw"
)

type CandidateViewWithDetails struct {
	AppID          uuid.UUID
	VacancyID      uuid.UUID
	CompanyID      uuid.UUID
	VacancyTitle   string
	CompanyName    string
	Status         application.Status
	CoverLetter    *string
	Snapshot       ApplicationSnapshot
	CreatedAt      time.Time
	UpdatedAt      time.Time
	StatusHistory  []StatusChangeCandidateView
	AllowedActions []AllowedAction
}

type StatusChangeCandidateView struct {
	Status        application.Status `json:"status"`
	ChangedByRole application.Actor  `json:"changed_by_role"`
	CreatedAt     time.Time          `json:"created_at"`
}

type ApplicationSnapshot struct {
	ResumeData projection.ResumeData
	Email      string
	FullName   string
	Telegram   *string
	CreatedAt  time.Time
}
