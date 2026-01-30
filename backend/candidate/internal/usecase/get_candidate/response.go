package get_candidate

import (
	"github.com/google/uuid"
	"time"
)

type CandidateResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Phone    string    `json:"phone"`
	Telegram string    `json:"telegram"`
	City     string    `json:"city"`
	Birthday time.Time `json:"birthday"`
}
