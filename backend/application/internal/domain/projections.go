package domain

import (
	"time"

	"github.com/google/uuid"
)

// VacancyStatus mirrors vacancy_status_enum.
type VacancyStatus string

const (
	VacancyStatusPublished VacancyStatus = "published"
	VacancyStatusArchived  VacancyStatus = "archived"
)

// CandidateProjection mirrors candidate_projection.
type CandidateProjection struct {
	ID uuid.UUID

	FullName string
	Email    string
	Telegram string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ResumeProjection mirrors resume_projection.
type ResumeProjection struct {
	ID uuid.UUID

	CandidateID uuid.UUID

	Name   string
	Data   ResumeData
	Status ResumeStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

// VacancyProjection mirrors vacancy_projection.
type VacancyProjection struct {
	ID uuid.UUID

	CompanyID   uuid.UUID
	CompanyName string

	Title  string
	Status VacancyStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

// CompanyMember mirrors company_members.
type CompanyMember struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Role      string
}
