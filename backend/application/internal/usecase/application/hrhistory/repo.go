package hrhistory

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type repo interface {
	GetHistoryHrView(ctx context.Context, appID, hrID uuid.UUID) ([]views.StatusChangeHrFullView, error)
}
