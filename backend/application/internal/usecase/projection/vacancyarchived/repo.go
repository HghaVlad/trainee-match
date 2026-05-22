package vacancyarchived

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type VacancyRepo interface {
	Archive(ctx context.Context, vacancyID uuid.UUID) error
}

type ApplicationRepo interface {
	UpdateStatusByVacancy(
		ctx context.Context,
		vacID uuid.UUID,
		newStatus application.Status,
		statusesToUpdate []application.Status,
		updAt time.Time,
	) error
}

type AppStatusHistoryRepo interface {
	AddChangesByVacancy(
		ctx context.Context,
		vacID uuid.UUID,
		newStatus application.Status,
		statusesToUpdate []application.Status,
		role application.Actor,
		comment *string,
		when time.Time,
	) error
}
