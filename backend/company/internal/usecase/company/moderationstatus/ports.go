package moderationstatus

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

type companyRepo interface {
	UpdateModerationStatusAndGetOld(
		ctx context.Context,
		companyID uuid.UUID,
		status company.ModerationStatus,
		when time.Time,
	) (company.ModerationStatus, error)
}

type outboxWriter interface {
	WriteCompanyModerationUpdated(ctx context.Context, ev company.ModerationUpdatedEvent) error
}

type cacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
