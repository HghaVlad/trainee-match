package candidatehistory

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Execute(
	ctx context.Context,
	appID uuid.UUID,
	ident identity.Identity,
) ([]views.StatusChangeCandidateFullView, error) {
	if ident.Role != identity.RoleCandidate {
		return nil, identity.ErrCandidateRoleRequired
	}

	history, err := u.repo.ListHistoryCandiView(ctx, appID, ident.UserID)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, application.ErrNotFound
	}

	return history, nil
}
