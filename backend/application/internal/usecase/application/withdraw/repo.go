package withdraw

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type appRepo interface {
	GetForUpdateByCandidate(
		ctx context.Context,
		appID, candID uuid.UUID,
	) (*application.Application, error)

	UpdateStatus(
		ctx context.Context,
		appID uuid.UUID,
		status application.Status,
		updAt time.Time,
	) error

	GetCandidateDetailedView(
		ctx context.Context,
		appID, candID uuid.UUID,
	) (*views.CandidateDetailedView, error)
}

type appStatusHistoryRepo interface {
	Add(ctx context.Context, change application.StatusChange) error
}
