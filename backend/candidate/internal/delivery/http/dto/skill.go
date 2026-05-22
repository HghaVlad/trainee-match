package dto

import "github.com/google/uuid"

type SkillResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SkillRequest struct {
	Name string `json:"name"`
}
