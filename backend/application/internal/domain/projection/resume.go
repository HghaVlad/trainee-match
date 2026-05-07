package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	ID          uuid.UUID
	CandidateID uuid.UUID
	Name        string
	Data        ResumeData
	Status      ResumeStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

var (
	ErrResumeNotFound = errors.New("resume projection not found")
)
