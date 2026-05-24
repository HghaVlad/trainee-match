package create

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

// Usecase creates vacancy in draft status
type Usecase struct {
	vacancyRepo VacancyRepo
	compRepo    CompanyRepo
	outboxWriter outboxWriter
	txManager   common.TxManager
}

func NewUsecase(
	vacancyRepo VacancyRepo,
	compRepo CompanyRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		vacancyRepo: vacancyRepo,
		compRepo:    compRepo,
		outboxWriter:  outboxWriter,
		txManager:   txManager,
	}
}

// Execute creates vacancy in draft status
func (u *Usecase) Execute(ctx context.Context, request *Request, ident *identity.Identity) (*Response, error) {
	vac := vacancyFromReq(request, ident)

	if err := vac.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	// only member of company can create vacancy
	comp, err := u.compRepo.GetByMember(ctx, request.CompanyID, ident.UserID)
	if err != nil {
		return nil, err
	}

	err = u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		err = u.vacancyRepo.Create(ctx, vac)
		if err != nil {
			return err
		}

		return u.createEvent(ctx, vac, comp)
	})

	return &Response{ID: vac.ID}, nil
}

func (u *Usecase) createEvent(ctx context.Context, vac *vacancy.Vacancy, comp *company.Company) error {
	ev := vacancy.DraftCreatedEvent{
		EventID:     uuid.New(),
		VacancyID:   vac.ID,
		Title:       vac.Title,
		CompanyID:   vac.CompanyID,
		CompanyName: comp.Name,
		OccurredAt:  time.Now().UTC(),
	}

	return u.outboxWriter.WriteVacancyDraftCreated(ctx, ev)
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
