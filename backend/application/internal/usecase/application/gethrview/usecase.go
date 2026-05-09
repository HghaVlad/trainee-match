package gethrview

import (
	"context"

	"github.com/google/uuid"

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
) (*views.HrDetailedView, error) {
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	view, err := u.repo.GetHrDetailedView(ctx, appID, ident.UserID)
	if err != nil {
		return nil, err
	}

	view.AllowedActions = views.HRAllowedActions(view.Status)
	return view, nil
}
