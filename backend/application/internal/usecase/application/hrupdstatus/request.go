package hrupdstatus

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type Request struct {
	AppID   uuid.UUID
	Status  application.Status
	Comment *string
}

func (r *Request) validate() error {
	if err := r.Status.IsValid(); err != nil {
		return err
	}

	if r.Comment != nil && len([]rune(*r.Comment)) > application.MaxCommentLength {
		return application.ErrStatusChangeCommentTooLong
	}

	return nil
}
