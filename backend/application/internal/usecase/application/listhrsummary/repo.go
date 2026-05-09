package listhrsummary

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
)

type appRepo interface {
	ListHrAppSummaries(
		ctx context.Context,
		hrID uuid.UUID,
		statuses []application.Status,
		companyID, vacancyID *uuid.UUID,
		createdFrom, createdTo *time.Time,
		order cursors.HrSummaryOrder,
		cursor any,
		limit int,
	) ([]views.HrAppSummary, error)
}

type memProjRepo interface {
	IsMember(ctx context.Context, userID uuid.UUID, companyID uuid.UUID) (bool, error)
}

type vacProjRepo interface {
	// CheckHrAccess returns companyID if success, projection.ErrVacancyNotFound otherwise
	CheckHrAccess(ctx context.Context, userID, vacancyID uuid.UUID) (uuid.UUID, error)
}
