package eventhandler

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Decoder interface {
	DecodeResumeUpsertedEvent(ctx context.Context, data []byte) (projection.ResumeUpsertedEvent, error)
}

type ResumeUpsertedUsecase interface {
	Execute(ctx context.Context, event projection.ResumeUpsertedEvent) error
}
