package handlers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listhrsummary"
	appviews "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/withdraw"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Deps struct {
	Apply               applyUC
	ListCandidateApps   listCandidateAppsUC
	GetCandidateViewUC  getCandidateViewUC
	GetCandiStatHistory listCandidateStatusHistoryUC
	Withdraw            withdrawUC
	ListHrApps          listHrAppsUC
	Logger              *slog.Logger
}

type applyUC interface {
	Execute(context.Context, apply.Request, identity.Identity) (*appviews.CandidateViewWithDetails, error)
}

type getCandidateViewUC interface {
	Execute(ctx context.Context, appID uuid.UUID, ident identity.Identity) (*appviews.CandidateViewWithDetails, error)
}

type listCandidateStatusHistoryUC interface {
	Execute(
		ctx context.Context,
		appID uuid.UUID,
		ident identity.Identity,
	) ([]appviews.StatusChangeCandidateFullView, error)
}

type listCandidateAppsUC interface {
	Execute(
		ctx context.Context,
		req listcandidatesummary.Request,
		ident identity.Identity,
	) (*listcandidatesummary.Response, error)
}

type withdrawUC interface {
	Execute(
		ctx context.Context,
		req withdraw.Request,
		ident identity.Identity,
	) (*appviews.CandidateViewWithDetails, error)
}

type listHrAppsUC interface {
	Execute(
		ctx context.Context,
		req listhrsummary.Request,
		ident identity.Identity,
	) (*listhrsummary.Response, error)
}
