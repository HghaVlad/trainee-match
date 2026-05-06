package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Candidate struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Telegram  *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

var (
	ErrCandidateNotFound = errors.New("candidate projection not found")
)
