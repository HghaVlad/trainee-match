package update_vacancy

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	repo      VacancyRepo
	cache     CacheRepo
	txManager uc_common.TxManager
}

func NewUsecase(
	repo VacancyRepo,
	cacheRepo CacheRepo,
	txManager uc_common.TxManager,
) *Usecase {

	return &Usecase{
		repo:      repo,
		cache:     cacheRepo,
		txManager: txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {

		vacancy, err := u.repo.GetByID(ctx, req.VacancyID, req.CompanyID)
		if err != nil {
			return err
		}

		applyPatch(vacancy, req)

		if vErr := vacancy.Validate(); vErr != nil {
			return vErr
		}

		return u.repo.Update(ctx, vacancy)
	})
	if err != nil {
		return err
	}

	u.cache.Del(ctx, req.VacancyID)
	return nil
}

// Applies not-nil only
func applyPatch(v *domain.Vacancy, r *Request) {
	if r.Title != nil {
		v.Title = *r.Title
	}

	if r.Description != nil {
		v.Description = *r.Description
	}

	if r.WorkFormat != nil {
		v.WorkFormat = *r.WorkFormat
	}

	if r.City != nil {
		v.City = r.City
	}

	if r.DurationFromMonths != nil {
		v.DurationFromMonths = r.DurationFromMonths
	}

	if r.DurationToMonths != nil {
		v.DurationToMonths = r.DurationToMonths
	}

	if r.EmploymentType != nil {
		v.EmploymentType = *r.EmploymentType
	}

	if r.HoursPerWeekFrom != nil {
		v.HoursPerWeekFrom = r.HoursPerWeekFrom
	}

	if r.HoursPerWeekTo != nil {
		v.HoursPerWeekTo = r.HoursPerWeekTo
	}

	if r.FlexibleSchedule != nil {
		v.FlexibleSchedule = *r.FlexibleSchedule
	}

	if r.IsPaid != nil {
		v.IsPaid = *r.IsPaid
	}

	if r.SalaryFrom != nil {
		v.SalaryFrom = r.SalaryFrom
	}

	if r.SalaryTo != nil {
		v.SalaryTo = r.SalaryTo
	}

	if r.InternshipToOffer != nil {
		v.InternshipToOffer = *r.InternshipToOffer
	}
}
