package create_vacancy

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	vacancyRepo VacancyRepo
	companyRepo CompanyRepo
	txManager   uc_common.TxManager
}

func NewUsecase(
	vacancyRepo VacancyRepo,
	companyRepo CompanyRepo,
	txManager uc_common.TxManager,
) *Usecase {

	return &Usecase{
		vacancyRepo: vacancyRepo,
		companyRepo: companyRepo,
		txManager:   txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, request *Request) (*Response, error) {
	vacancy := vacancyFromReq(request)

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {

		vacErr := u.vacancyRepo.Create(ctx, vacancy)
		if vacErr != nil {
			return vacErr
		}

		return u.companyRepo.IncrementOpenVacancies(ctx, vacancy.CompanyID)
	})

	if err != nil {
		return nil, err
	}

	return &Response{ID: vacancy.ID}, nil
}

func vacancyFromReq(request *Request) *domain.Vacancy {
	return &domain.Vacancy{
		ID:        uuid.New(),
		CompanyID: request.CompanyID,

		Title:       request.Title,
		Description: request.Description,

		WorkFormat: request.WorkFormat,
		City:       request.City,

		DurationFromMonths: request.DurationFromMonths,
		DurationToMonths:   request.DurationToMonths,

		EmploymentType:   request.EmploymentType,
		HoursPerWeekFrom: request.HoursPerWeekFrom,
		HoursPerWeekTo:   request.HoursPerWeekTo,

		FlexibleSchedule: request.FlexibleSchedule,

		IsPaid:     request.IsPaid,
		SalaryFrom: request.SalaryFrom,
		SalaryTo:   request.SalaryTo,

		InternshipToOffer: request.InternshipToOffer,
	}
}
