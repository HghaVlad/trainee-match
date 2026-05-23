package companydeleted

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type CompanyMemberRepo interface {
	DeleteByCompanyID(ctx context.Context, companyID uuid.UUID) error
}

type VacancyRepo interface {
	ArchiveByCompany(ctx context.Context, compID uuid.UUID) error
}

type ApplicationRepo interface {
	UpdateStatusByCompany(
		ctx context.Context,
		compID uuid.UUID,
		newStatus application.Status,
		statusesToUpdate []application.Status,
		updAt time.Time,
	) error
}

type AppStatusHistoryRepo interface {
	AddChangesByCompany(
		ctx context.Context,
		compID uuid.UUID,
		newStatus application.Status,
		statusesToUpdate []application.Status,
		role application.Actor,
		comment *string,
		when time.Time,
	) error
}
