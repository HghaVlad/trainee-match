package getcandidateview

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	repo repo
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Execute(
	ctx context.Context,
	appID uuid.UUID,
	ident identity.Identity,
) (*views.CandidateViewWithDetails, error) {
	if ident.Role != identity.RoleCandidate {
		return nil, identity.ErrCandidateRoleRequired
	}

	view, err := u.repo.GetByIDCandidateViewWithDetails(ctx, appID, ident.UserID)
	if err != nil {
		return nil, err
	}

	if view.Status != application.StatusRejected && view.Status != application.StatusOffer &&
		view.Status != application.StatusWithdrawn {
		view.AllowedActions = append(view.AllowedActions, views.AllowedActionWithdraw)
	}

	return view, nil
}
