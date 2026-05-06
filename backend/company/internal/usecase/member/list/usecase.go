package list

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	repo repo
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req Request, ident *identity.Identity) (*Response, error) {
	if err := req.IsValid(); err != nil {
		return nil, err
	}

	if err := u.authorize(ctx, req.CompanyID, ident); err != nil {
		return nil, err
	}

	views, err := u.repo.ListViewsByCompany(ctx, req.CompanyID, req.Limit+1, req.Offset)
	if err != nil {
		return nil, err
	}

	resp := &Response{}

	if len(views) == req.Limit+1 {
		resp.HasMore = true
		resp.Members = views[:req.Limit]
	} else {
		resp.HasMore = false
		resp.Members = views
	}

	return resp, nil
}

// only member of company can archive vacancy
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, iden *identity.Identity) error {
	if iden.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	_, err := u.repo.Get(ctx, iden.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}

	return err
}
