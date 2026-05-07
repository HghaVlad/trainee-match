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
		statuses []application.Status,
		companyID *uuid.UUID,
		vacancyID *uuid.UUID,
		createdFrom *time.Time,
		createdTo *time.Time,
		order cursors.HrSummaryOrder,
		cursor any,
		limit int,
	) ([]views.HrAppSummary, error)
}

type memProjRepo interface {
	IsMember(ctx context.Context, userID uuid.UUID, companyID uuid.UUID) (bool, error)
}

type vacProjRepo interface {
	GetCompanyIDByVacancyID(ctx context.Context, vacancyID uuid.UUID) (uuid.UUID, error)
}
