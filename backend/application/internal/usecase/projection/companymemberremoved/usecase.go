package companymemberremoved

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

func (uc *Usecase) Execute(ctx context.Context, event projection.CompanyMemberRemovedEvent) error {
	return uc.repo.Delete(ctx, event.UserID, event.CompanyID)
}
