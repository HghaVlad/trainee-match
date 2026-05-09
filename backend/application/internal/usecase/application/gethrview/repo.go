package gethrview

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type repo interface {
	// GetHrDetailedView returns detailed view,
	// if hr is not a member of company owning this application, returns application.ErrNotFound
	GetHrDetailedView(
		ctx context.Context,
		appID, hrID uuid.UUID,
	) (*views.HrDetailedView, error)
}
