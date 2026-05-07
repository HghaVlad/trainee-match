package getcandidateview

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type repo interface {
	GetByIDCandidateViewWithDetails(
		ctx context.Context,
		appID, candID uuid.UUID,
	) (*views.CandidateViewWithDetails, error)
}
