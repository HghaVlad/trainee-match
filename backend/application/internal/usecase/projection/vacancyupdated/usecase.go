package vacancyupdated

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	repo VacancyRepo
}

func NewUsecase(repo VacancyRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) Execute(ctx context.Context, event projection.VacancyUpdatedEvent) error {
	return uc.repo.UpdateTitle(ctx, event.VacancyID, event.Title)
}
