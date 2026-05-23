package companydeleted

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common"
)

type Usecase struct {
	memberRepo    CompanyMemberRepo
	vacancyRepo   VacancyRepo
	appRepo       ApplicationRepo
	appStatusRepo AppStatusHistoryRepo
	txManager     common.TxManager
}

func NewUsecase(
	memberRepo CompanyMemberRepo,
	vacancyRepo VacancyRepo,
	appRepo ApplicationRepo,
	appStatusRepo AppStatusHistoryRepo,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		memberRepo:    memberRepo,
		vacancyRepo:   vacancyRepo,
		appRepo:       appRepo,
		appStatusRepo: appStatusRepo,
		txManager:     txManager,
	}
}

// Execute removes company members, archives company's vacancies,
// active applications to company's vacancies are rejected
func (uc *Usecase) Execute(ctx context.Context, event projection.CompanyDeletedEvent) error {
	now := time.Now().UTC()

	comment := "Application rejected because company was deleted"

	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.memberRepo.DeleteByCompanyID(ctx, event.CompanyID); err != nil {
			return err
		}

		if err := uc.vacancyRepo.ArchiveByCompany(ctx, event.CompanyID); err != nil {
			return err
		}

		err := uc.appStatusRepo.AddChangesByCompany(
			ctx,
			event.CompanyID,
			application.StatusRejected,
			application.ActiveStatuses(),
			application.ActorSystem,
			&comment,
			now,
		)
		if err != nil {
			return err
		}

		return uc.appRepo.UpdateStatusByCompany(
			ctx,
			event.CompanyID,
			application.StatusRejected,
			application.ActiveStatuses(),
			now,
		)
	})
}
