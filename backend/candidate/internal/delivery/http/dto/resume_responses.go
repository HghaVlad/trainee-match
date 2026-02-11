package dto

import (
	"github.com/google/uuid"
)

type ResumeResponse struct {
	ID          uuid.UUID  `json:"id"`
	CandidateId uuid.UUID  `json:"candidate_id"`
	Name        string     `json:"name"`
	Status      int        `json:"status"`
	Data        ResumeData `json:"data"`
}

type ShortResumeResponse struct {
	ID          uuid.UUID `json:"id"`
	CandidateId uuid.UUID `json:"candidate_id"`
	Name        string    `json:"name"`
	Status      int       `json:"status"`
}

type SkillResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
