package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

// Usecase creates vacancy in draft status
type Usecase struct {
	vacancyRepo VacancyRepo
	memberRepo  CompMemberRepo
	compRepo    CompanyRepo
	searchRepo  SearchRepo
}

func NewUsecase(
	vacancyRepo VacancyRepo,
	memberRepo CompMemberRepo,
	compRepo CompanyRepo,
	searchRepo SearchRepo,
) *Usecase {
	return &Usecase{
		vacancyRepo: vacancyRepo,
		memberRepo:  memberRepo,
		compRepo:    compRepo,
		searchRepo:  searchRepo,
	}
}

// TODO: add everywhere where vacancy is updated and company name is updated

// Execute creates vacancy in draft status
func (u *Usecase) Execute(ctx context.Context, request *Request, ident *identity.Identity) (*Response, error) {
	vac := vacancyFromReq(request, ident)

	if err := vac.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// only member of company can create vacancy
	comp, err := u.compRepo.GetByMember(ctx, request.CompanyID, ident.UserID)
	if err != nil {
		return nil, err
	}

	err = u.vacancyRepo.Create(ctx, vac)
	if err != nil {
		return nil, err
	}

	err = u.searchRepo.Index(ctx, *vac, comp.Name)
	if err != nil {
		return nil, err
	}

	return &Response{ID: vac.ID}, nil
}

// user of identity is the creator of the vacancy
func vacancyFromReq(request *Request, ident *identity.Identity) *vacancy.Vacancy {
	vac := &vacancy.Vacancy{
		ID:        uuid.New(),
		CompanyID: request.CompanyID,
		CreatedBy: ident.UserID,

		Title:       request.Title,
		Description: request.Description,

		Status: vacancy.StatusDraft,

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

		ModerationStatus: vacancy.ModerationStatusOK,
		CreatedAt:        time.Now().UTC(),
	}

	if request.EmploymentType != nil {
		vac.EmploymentType = *request.EmploymentType
	} else {
		vac.EmploymentType = vacancy.EmploymentTypeInternship
	}

	return vac
}
