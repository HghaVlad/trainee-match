package create_vacancy

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
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
	id := uuid.New()
	vacancy := vacancyFromReq(request, id)

	dErr := vacancy.Validate()
	if dErr != nil {
		return nil, dErr
	}

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

func vacancyFromReq(request *Request, id uuid.UUID) *domain.Vacancy {
	vacancy := &domain.Vacancy{
		ID:        id,
		CompanyID: request.CompanyID,

		Title:       request.Title,
		Description: request.Description,

		WorkFormat: request.WorkFormat,
		City:       request.City,

		DurationFromMonths: request.DurationFromMonths,
		DurationToMonths:   request.DurationToMonths,

		HoursPerWeekFrom: request.HoursPerWeekFrom,
		HoursPerWeekTo:   request.HoursPerWeekTo,

		FlexibleSchedule: request.FlexibleSchedule,

		IsPaid:     request.IsPaid,
		SalaryFrom: request.SalaryFrom,
		SalaryTo:   request.SalaryTo,

		InternshipToOffer: request.InternshipToOffer,
	}

	if request.EmploymentType != nil {
		vacancy.EmploymentType = *request.EmploymentType
	} else {
		vacancy.EmploymentType = value_types.EmploymentTypeInternship
	}

	return vacancy
}
