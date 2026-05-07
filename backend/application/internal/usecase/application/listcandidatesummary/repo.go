package listcandidatesummary

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type repo interface {
	ListCandidateSummaries(
		ctx context.Context,
		candidateID uuid.UUID,
		statuses []application.Status,
		companyID *uuid.UUID,
		order Order,
		cursor any,
		limit int,
	) ([]views.CandidateSummary, error)
}
