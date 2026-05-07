package hash

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type AppSnapshotHasher struct{}

func NewAppSnapshotHasher() *AppSnapshotHasher {
	return &AppSnapshotHasher{}
}

type snapshotPayload struct {
	ResumeData projection.ResumeData `json:"resume_data"`

	CandidateFullName string  `json:"candidate_full_name"`
	CandidateEmail    string  `json:"candidate_email"`
	CandidateTelegram *string `json:"candidate_telegram,omitempty"`
}

func (h *AppSnapshotHasher) GetDeterministicAppSnapshotID(
	resume projection.ResumeData,
	candidate projection.Candidate,
) (uuid.UUID, error) {
	truncateDateOfBirth(&resume)

	payload := snapshotPayload{
		ResumeData: resume,

		CandidateFullName: candidate.FullName,
		CandidateEmail:    candidate.Email,
		CandidateTelegram: candidate.Telegram,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("marshal snapshot payload: %w", err)
	}

	sum := sha256.Sum256(data)

	var id uuid.UUID
	copy(id[:], sum[:16])
	id[6] = (id[6] & 0x0f) | 0x50 // version 5
	id[8] = (id[8] & 0x3f) | 0x80 // RFC 4122 variant

	return id, nil
}

func truncateDateOfBirth(resume *projection.ResumeData) {
	resume.DateOfBirth = time.Date(
		resume.DateOfBirth.Year(),
		resume.DateOfBirth.Month(),
		resume.DateOfBirth.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
}
