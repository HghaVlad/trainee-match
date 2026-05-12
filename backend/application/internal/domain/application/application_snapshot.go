package application

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Snapshot struct {
	ID          uuid.UUID
	ResumeID    uuid.UUID
	CandidateID uuid.UUID

	ResumeName string
	ResumeData projection.ResumeData

	FullName string
	Email    string
	Telegram *string

	CreatedAt time.Time
}

func NewApplicationSnapshot(
	id uuid.UUID,
	res projection.Resume,
	cand projection.Candidate,
	createdAt time.Time,
) (*Snapshot, error) {
	return &Snapshot{
		ID:          id,
		ResumeID:    res.ID,
		CandidateID: cand.ID,
		ResumeName:  res.Name,
		ResumeData:  res.Data,
		FullName:    cand.FullName,
		Email:       cand.Email,
		Telegram:    cand.Telegram,
		CreatedAt:   createdAt,
	}, nil
}
