package dto

import (
	"github.com/google/uuid"
	"time"
)

type ResumeResponse struct {
	ID          uuid.UUID  `json:"id"`
	CandidateId uuid.UUID  `json:"candidate_id"`
	Name        string     `json:"name"`
	Status      int        `json:"status"`
	Data        ResumeData `json:"data"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SkillResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

