package eventhandler

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

//go:generate mockgen -source=ports.go -destination=mocks/mocks.go -package=mocks
type DLQSender interface {
	ToDLQ(ctx context.Context, eventID uuid.UUID, key, payload []byte, topic, eventType string, errMsg string) error
}

type Decoder interface {
	GetUserCreatedEvent(ctx context.Context, payload []byte) (*userhr.CreatedEvent, error)
	GetVacancyPublishedEvent(ctx context.Context, payload []byte) (*vacancy.PublishedEvent, error)
	GetVacancyDraftCreatedEvent(ctx context.Context, payload []byte) (*vacancy.DraftCreatedEvent, error)
	GetVacancyUpdatedEvent(ctx context.Context, payload []byte) (*vacancy.UpdatedEvent, error)
	GetVacancyArchivedEvent(ctx context.Context, payload []byte) (*vacancy.ArchivedEvent, error)
	GetVacancyModUpdEvent(ctx context.Context, payload []byte) (*vacancy.ModerationUpdatedEvent, error)
	GetCompanyUpdatedEvent(ctx context.Context, payload []byte) (*company.UpdatedEvent, error)
	GetCompanyDeletedEvent(ctx context.Context, payload []byte) (*company.DeletedEvent, error)
}

type UserHrCreator interface {
	Execute(ctx context.Context, ev userhr.CreatedEvent) error
}

type SearchIndexer interface {
	Index(ctx context.Context, vacID uuid.UUID) error
	UpdateCompanyName(ctx context.Context, compID uuid.UUID, name string) error
	RemoveByCompany(ctx context.Context, compID uuid.UUID) error
}
