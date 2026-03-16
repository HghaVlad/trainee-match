package create_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

// Usecase creates vacancy in draft status
type Usecase struct {
	vacancyRepo VacancyRepo
	memberRepo  CompMemberRepo
}

func NewUsecase(
	vacancyRepo VacancyRepo,
	memberRepo CompMemberRepo,
) *Usecase {

	return &Usecase{
		vacancyRepo: vacancyRepo,
		memberRepo:  memberRepo,
	}
}

// Execute creates vacancy in draft status
func (u *Usecase) Execute(ctx context.Context, request *Request, identity uc_common.Identity) (*Response, error) {
	vacancy := vacancyFromReq(request, identity)

	if err := vacancy.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := u.authorize(ctx, request.CompanyID, identity); err != nil {
		return nil, err
	}

	err := u.vacancyRepo.Create(ctx, vacancy)
	if err != nil {
		return nil, err
	}

	return &Response{ID: vacancy.ID}, nil
}

// only member of company can create vacancy
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, identity uc_common.Identity) error {
	if identity.Role != uc_common.RoleHR {
		return domain_errors.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, identity.UserID, companyID)
	if errors.Is(err, domain_errors.ErrCompanyMemberNotFound) {
		return domain_errors.ErrCompanyMemberRequired
	}

	return err
}

// user of identity is the creator of the vacancy
func vacancyFromReq(request *Request, identity uc_common.Identity) *domain.Vacancy {
	vacancy := &domain.Vacancy{
		ID:        uuid.New(),
		CompanyID: request.CompanyID,
		CreatedBy: identity.UserID,

		Title:       request.Title,
		Description: request.Description,

		Status: value_types.VacancyStatusDraft,

		WorkFormat: request.WorkFormat,
		City:       request.City,

		DurationFromDays: request.DurationFromDays,
		DurationToDays:   request.DurationToDays,

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
