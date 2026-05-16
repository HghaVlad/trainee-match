package companydeleted

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	memberRepo  CompanyMemberRepo
	vacancyRepo VacancyRepo
}

func NewUsecase(memberRepo CompanyMemberRepo, vacancyRepo VacancyRepo) *Usecase {
	return &Usecase{
		memberRepo:  memberRepo,
		vacancyRepo: vacancyRepo,
	}
}

func (uc *Usecase) Execute(ctx context.Context, event projection.CompanyDeletedEvent) error {
	if err := uc.memberRepo.DeleteByCompanyID(ctx, event.CompanyID); err != nil {
		return err
	}
	return uc.vacancyRepo.DeleteByCompanyID(ctx, event.CompanyID)
}
