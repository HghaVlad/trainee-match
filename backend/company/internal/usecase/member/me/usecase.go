package me

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/views"
)

type Usecase struct {
	repo repo
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(
	ctx context.Context,
	compID uuid.UUID,
	iden identity.Identity,
) (*views.MemberFullView, error) {
	if iden.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	mem, err := u.repo.GetFullView(ctx, iden.UserID, compID)
	if err != nil {
		return nil, err
	}

	return mem, nil
}
