package listmy

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
)

type Usecase struct {
	compList *list.Usecase
}

func NewUsecase(compList *list.Usecase) *Usecase {
	return &Usecase{
		compList: compList,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, ident *identity.Identity) (*list.Response, error) {
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	listCompReq := list.Request{
		Order:         req.Order,
		Limit:         req.Limit,
		EncodedCursor: req.EncodedCursor,
		Filter: list.Filter{
			CompanyMemberID: &ident.UserID,
		},
	}

	return u.compList.Execute(ctx, &listCompReq)
}
