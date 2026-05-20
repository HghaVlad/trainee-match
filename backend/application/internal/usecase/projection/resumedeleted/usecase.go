package resumedeleted

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	repo ResumeRepo
}

func NewUsecase(repo ResumeRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) Execute(ctx context.Context, event projection.ResumeDeletedEvent) error {
	return uc.repo.Delete(ctx, event.ResumeID)
}
