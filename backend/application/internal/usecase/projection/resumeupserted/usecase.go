package resumeupserted

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

func (uc *Usecase) Execute(ctx context.Context, event projection.ResumeUpsertedEvent) error {
	resume := event.ToResume()
	err := uc.repo.Save(ctx, &resume)
	return err
}
