package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrCandidateNotFound = errors.New("candidate not found")
)

type Candidate struct {
	ID       uuid.UUID
	UserId   uuid.UUID
	Phone    string
	Telegram string
	City     string
	Birthday time.Time
}
