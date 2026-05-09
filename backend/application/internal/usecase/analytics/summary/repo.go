package summary

import (
	"context"

	"github.com/google/uuid"
)

type appRepo interface {
	GetCompanyAnalyticsSummary(ctx context.Context, compID uuid.UUID) (*Summary, error)
	GetVacancyAnalyticsSummary(ctx context.Context, vacID uuid.UUID) (*Summary, error)
}

type companyMemberRepo interface {
	IsMember(
		ctx context.Context,
		userID uuid.UUID,
		companyID uuid.UUID,
	) (bool, error)
}

type vacancyProjRepo interface {
	CheckHrAccess(
		ctx context.Context,
		userID, vacancyID uuid.UUID,
	) (uuid.UUID, error)
}
