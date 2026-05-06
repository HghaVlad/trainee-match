package apply

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type candidateProjRepo interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*projection.Candidate, error)
}

type resumeProjRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*projection.Resume, error)
}

type vacProjRepo interface {
	GetByID(ctx context.Context, vacID uuid.UUID) (*projection.Vacancy, error)
}

type appSnapshotRepo interface {
	CreateIdempotent(ctx context.Context, appSnapshot application.Snapshot) error
}

type appRepo interface {
	Create(ctx context.Context, app application.Application) error
}

type appHistoryRepo interface {
	Add(ctx context.Context, change application.StatusChange) error
}
