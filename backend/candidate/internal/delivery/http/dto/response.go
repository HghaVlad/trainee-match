package dto

import (
	"github.com/google/uuid"
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type CandidateResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Phone    string    `json:"phone"`
	Telegram string    `json:"telegram"`
	City     string    `json:"city"`
	Birthday time.Time `json:"birthday"`
}
