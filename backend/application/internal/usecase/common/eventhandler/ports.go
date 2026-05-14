package eventhandler

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Decoder interface {
	DecodeResumeUpsertedEvent(ctx context.Context, data []byte) (projection.ResumeUpsertedEvent, error)
	DecodeResumeDeletedEvent(ctx context.Context, data []byte) (projection.ResumeDeletedEvent, error)
	DecodeCandidateUpsertedEvent(ctx context.Context, data []byte) (projection.CandidateUpsertedEvent, error)
}

type ResumeUpsertedUsecase interface {
	Execute(ctx context.Context, event projection.ResumeUpsertedEvent) error
}

type ResumeDeletedUsecase interface {
	Execute(ctx context.Context, event projection.ResumeDeletedEvent) error
}

type CandidateUpsertedUsecase interface {
	Execute(ctx context.Context, event projection.CandidateUpsertedEvent) error
}
