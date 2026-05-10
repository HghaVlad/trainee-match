package withdraw

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type Request struct {
	AppID   uuid.UUID
	Comment *string
}

func (r *Request) validate() error {
	if r.Comment != nil && len([]rune(*r.Comment)) > application.MaxCommentLength {
		return application.ErrStatusChangeCommentTooLong
	}

	return nil
}
