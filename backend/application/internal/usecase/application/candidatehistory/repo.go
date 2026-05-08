package candidatehistory

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type Repo interface {
	ListHistoryCandiView(ctx context.Context, appID, candID uuid.UUID) ([]views.StatusChangeCandidateFullView, error)
}
