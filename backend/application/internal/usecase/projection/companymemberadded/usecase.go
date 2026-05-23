package companymemberadded

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	repo CompanyMemberRepo
}

func NewUsecase(repo CompanyMemberRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) Execute(ctx context.Context, event projection.CompanyMemberAddedEvent) error {
	member := event.ToCompanyMember()
	return uc.repo.Save(ctx, &member)
}
