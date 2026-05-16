package vacancypublished

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

func (uc *Usecase) Execute(ctx context.Context, event projection.VacancyPublishedEvent) error {
	vacancy := event.ToVacancy()
	return uc.repo.Save(ctx, &vacancy)
}
