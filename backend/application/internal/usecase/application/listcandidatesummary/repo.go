package listcandidatesummary

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
)

type repo interface {
	ListCandidateAppSummaries(
		ctx context.Context,
		candidateID uuid.UUID,
		statuses []application.Status,
		companyID *uuid.UUID,
		order cursors.SummaryOrder,
		cursor any,
		limit int,
	) ([]views.CandidateAppSummary, error)
}
