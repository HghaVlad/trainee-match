package apply

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type Request struct {
	VacancyID   uuid.UUID
	ResumeID    uuid.UUID
	CoverLetter *string
}

func (r *Request) validate() error {
	if r.CoverLetter != nil && len([]rune(*r.CoverLetter)) > application.MaxCoverLetterLength {
		return application.ErrCoverLetterTooLong
	}

	return nil
}
