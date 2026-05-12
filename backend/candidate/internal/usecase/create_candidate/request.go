package create_candidate

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Telegram string    `json:"telegram"`
	City     string    `json:"city"`
	Birthday time.Time `json:"birthday"`
}
