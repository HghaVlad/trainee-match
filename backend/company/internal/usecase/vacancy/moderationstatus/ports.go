package moderationstatus

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type UpdateModerationResult struct {
	CompanyID           uuid.UUID
	OldModerationStatus vacancy.ModerationStatus
	VacancyStatus       vacancy.Status
}

//go:generate mockgen -source=ports.go -destination=mocks/port_mocks.go -package=mocks
type vacancyRepo interface {
	UpdateModerationStatus(
		ctx context.Context,
		vacancyID uuid.UUID,
		status vacancy.ModerationStatus,
		when time.Time,
	) (*UpdateModerationResult, error)
}

type companyRepo interface {
	DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error
	IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error
}

type outboxWriter interface {
	WriteVacancyModerationUpdated(ctx context.Context, ev vacancy.ModerationUpdatedEvent) error
}

type cacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
