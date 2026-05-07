package handlers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
	appviews "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Deps struct {
	Apply              applyUC
	ListCandidateApps  listCandidateAppsUC
	GetCandidateViewUC getCandidateViewUC
	Logger             *slog.Logger
}

type applyUC interface {
	Execute(context.Context, apply.Request, identity.Identity) (*appviews.CandidateViewWithDetails, error)
}

type getCandidateViewUC interface {
	Execute(ctx context.Context, appID uuid.UUID, ident identity.Identity) (*appviews.CandidateViewWithDetails, error)
}

type listCandidateAppsUC interface {
	Execute(
		ctx context.Context,
		req listcandidatesummary.Request,
		ident identity.Identity,
	) (*listcandidatesummary.Response, error)
}
