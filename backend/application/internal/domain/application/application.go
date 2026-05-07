package application

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusSubmitted Status = "submitted"
	StatusSeen      Status = "seen"
	StatusInterview Status = "interview"
	StatusRejected  Status = "rejected"
	StatusOffer     Status = "offer"
	StatusWithdrawn Status = "withdrawn"
)

type Application struct {
	ID          uuid.UUID
	ResumeID    uuid.UUID
	CandidateID uuid.UUID
	VacancyID   uuid.UUID
	CompanyID   uuid.UUID
	SnapshotID  uuid.UUID
	Status      Status
	CoverLetter *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSubmitted(
	resumeID, candID, vacID, compID, snapID uuid.UUID,
	coverLetter *string,
	when time.Time,
) (*Application, error) {
	return &Application{
		ID:          uuid.New(),
		ResumeID:    resumeID,
		CandidateID: candID,
		VacancyID:   vacID,
		CompanyID:   compID,
		SnapshotID:  snapID,
		Status:      StatusSubmitted,
		CoverLetter: coverLetter,
		CreatedAt:   when,
		UpdatedAt:   when,
	}, nil
}
