package getcandidate

import (
	"time"

	"github.com/google/uuid"
)

type CandidateResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	UserID   uuid.UUID `json:"user_id"`
	Phone    string    `json:"phone"`
	Telegram string    `json:"telegram"`
	City     string    `json:"city"`
	Birthday time.Time `json:"birthday"`
}
