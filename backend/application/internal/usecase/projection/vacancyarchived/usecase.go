package vacancyarchived

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common"
)

type Usecase struct {
	repo          VacancyRepo
	appRepo       ApplicationRepo
	appStatusRepo AppStatusHistoryRepo
	txManager     common.TxManager
}

func NewUsecase(
	repo VacancyRepo,
	appRepo ApplicationRepo,
	appStatus AppStatusHistoryRepo,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		repo:          repo,
		appRepo:       appRepo,
		appStatusRepo: appStatus,
		txManager:     txManager,
	}
}

// Execute archives vacancy, all active applications to this vacancy become rejected
func (uc *Usecase) Execute(ctx context.Context, event projection.VacancyArchivedEvent) error {
	now := time.Now().UTC()

	comment := "Application rejected because vacancy was archived"

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		err := uc.repo.Archive(ctx, event.VacancyID)
		if err != nil {
			return err
		}

		err = uc.appStatusRepo.AddChangesByVacancy(
			ctx,
			event.VacancyID,
			application.StatusRejected,
			application.ActiveStatuses(),
			application.ActorSystem,
			&comment,
			now,
		)
		if err != nil {
			return err
		}

		return uc.appRepo.UpdateStatusByVacancy(
			ctx,
			event.VacancyID,
			application.StatusRejected,
			application.ActiveStatuses(),
			now,
		)
	})
}
