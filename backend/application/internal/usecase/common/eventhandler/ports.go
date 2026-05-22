package eventhandler

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/dlq"
)

type Decoder interface {
	DecodeResumeUpsertedEvent(ctx context.Context, data []byte) (projection.ResumeUpsertedEvent, error)
	DecodeResumeDeletedEvent(ctx context.Context, data []byte) (projection.ResumeDeletedEvent, error)
	DecodeCandidateUpsertedEvent(ctx context.Context, data []byte) (projection.CandidateUpsertedEvent, error)
	DecodeCompanyUpdatedEvent(ctx context.Context, data []byte) (projection.CompanyUpdatedEvent, error)
	DecodeCompanyDeletedEvent(ctx context.Context, data []byte) (projection.CompanyDeletedEvent, error)
	DecodeCompanyMemberAddedEvent(ctx context.Context, data []byte) (projection.CompanyMemberAddedEvent, error)
	DecodeCompanyMemberRemovedEvent(ctx context.Context, data []byte) (projection.CompanyMemberRemovedEvent, error)
	DecodeVacancyPublishedEvent(ctx context.Context, data []byte) (projection.VacancyPublishedEvent, error)
	DecodeVacancyArchivedEvent(ctx context.Context, data []byte) (projection.VacancyArchivedEvent, error)
	DecodeVacancyUpdatedEvent(ctx context.Context, data []byte) (projection.VacancyUpdatedEvent, error)
	DecodeVacancyModerationUpdEvent(ctx context.Context, data []byte) (projection.VacancyModerationUpdatedEvent, error)
	DecodeCompanyModerationUpdEvent(ctx context.Context, data []byte) (projection.CompanyModerationUpdatedEvent, error)
}

type DLQSender interface {
	SendDLQ(ctx context.Context, message dlq.Message, key []byte) error
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

type CompanyUpdatedUsecase interface {
	Execute(ctx context.Context, event projection.CompanyUpdatedEvent) error
}

type CompanyDeletedUsecase interface {
	Execute(ctx context.Context, event projection.CompanyDeletedEvent) error
}

type CompanyMemberAddedUsecase interface {
	Execute(ctx context.Context, event projection.CompanyMemberAddedEvent) error
}

type CompanyMemberRemovedUsecase interface {
	Execute(ctx context.Context, event projection.CompanyMemberRemovedEvent) error
}

type VacancyPublishedUsecase interface {
	Execute(ctx context.Context, event projection.VacancyPublishedEvent) error
}

type VacancyArchivedUsecase interface {
	Execute(ctx context.Context, event projection.VacancyArchivedEvent) error
}

type VacancyUpdatedUsecase interface {
	Execute(ctx context.Context, event projection.VacancyUpdatedEvent) error
}

type VacancyModerationUpdUsecase interface {
	Execute(ctx context.Context, event projection.VacancyModerationUpdatedEvent) error
}

type CompanyModerationUpdUsecase interface {
	Execute(ctx context.Context, event projection.CompanyModerationUpdatedEvent) error
}
