package vacancymodupd

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	repo vacancyRepo
}

func NewUsecase(repo vacancyRepo) Usecase {
	return Usecase{
		repo: repo,
	}
}

func (u Usecase) Execute(ctx context.Context, event projection.VacancyModerationUpdatedEvent) error {
	return u.repo.UpdateModStatus(ctx, event.VacancyID, event.ModStatus)
}
