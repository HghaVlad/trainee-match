package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	ApplicationStatusSubmitted ApplicationStatus = "submitted"
	ApplicationStatusSeen      ApplicationStatus = "seen"
	ApplicationStatusInterview ApplicationStatus = "interview"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
	ApplicationStatusOffer     ApplicationStatus = "offer"
	ApplicationStatusWithdrawn ApplicationStatus = "withdrawn"
)

type ApplicationSnapshot struct {
	ID          uuid.UUID
	ResumeID    uuid.UUID
	CandidateID uuid.UUID

	ResumeName string
	ResumeData ResumeData

	FullName string
	Email    string
	Telegram string

	Hash      string
	CreatedAt time.Time
}

// Application mirrors applications.
type Application struct {
	ID          uuid.UUID
	ResumeID    uuid.UUID
	CandidateID uuid.UUID
	VacancyID   uuid.UUID
	CompanyID   uuid.UUID

	SnapshotID uuid.UUID

	Status      ApplicationStatus
	CoverLetter string

	CreatedAt time.Time
	UpdatedAt time.Time
}
