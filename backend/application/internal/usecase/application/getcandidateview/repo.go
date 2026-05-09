package getcandidateview

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type repo interface {
	// GetCandidateDetailedView returns detailed view,
	// if candidate doesn't own this application, returns application.ErrNotFound
	GetCandidateDetailedView(
		ctx context.Context,
		appID, candID uuid.UUID,
	) (*views.CandidateDetailedView, error)
}
