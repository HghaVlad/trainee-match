package update

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type Usecase struct {
	repo       VacancyRepo
	compRepo   compRepo
	outbox     outboxWriter
	searchRepo searchRepo
	cache      CacheRepo
	txManager  common.TxManager
}

func NewUsecase(
	repo VacancyRepo,
	compRepo compRepo,
	outbox outboxWriter,
	searchRepo searchRepo,
	cacheRepo CacheRepo,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		repo:       repo,
		compRepo:   compRepo,
		outbox:     outbox,
		searchRepo: searchRepo,
		cache:      cacheRepo,
		txManager:  txManager,
	}
}

// Execute updates vacancy. All nil fields of vacancy in request won't be applied.
// Deletes vacancy from cache.
func (u *Usecase) Execute(ctx context.Context, req *Request, ident *identity.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if err := req.lightValidate(); err != nil {
		return err
	}

	// only member of company can update vacancy
	comp, err := u.compRepo.GetByMember(ctx, req.CompanyID, ident.UserID)
	if err != nil {
		return err
	}

	var vacncy *vacancy.Vacancy

	err = u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		vac, err := u.repo.GetByIDForUpdate(ctx, req.VacancyID, req.CompanyID)
		if err != nil {
			return err
		}

		eventShouldBeCreated := checkIfEventShouldBeCreated(vac, req)

		applyPatch(vac, req)

		if vErr := vac.Validate(); vErr != nil {
			return vErr
		}

		err = u.repo.Update(ctx, vac)
		if err != nil {
			return err
		}

		vacncy = vac
		if eventShouldBeCreated {
			return u.createdVacancyUpdatedEvent(ctx, vac)
		}

		return nil
	})
	if err != nil {
		return err
	}

	searchView := views.SearchViewFromVacancy(*vacncy, comp.Name)
	if err := u.searchRepo.Index(ctx, *searchView); err != nil {
		return err
	}

	u.cache.Del(ctx, req.VacancyID)
	return nil
}

func (u *Usecase) createdVacancyUpdatedEvent(ctx context.Context, vac *vacancy.Vacancy) error {
	ev := vacancy.UpdatedEvent{
		EventID:    uuid.New(),
		VacancyID:  vac.ID,
		Title:      vac.Title,
		OccurredAt: time.Now().UTC(),
	}

	return u.outbox.WriteVacancyUpdated(ctx, ev)
}

func checkIfEventShouldBeCreated(vac *vacancy.Vacancy, req *Request) bool {
	return vac.Status != vacancy.StatusDraft && req.Title != nil && vac.Title != *req.Title
}

// Applies not-nil only
func applyPatch(v *vacancy.Vacancy, r *Request) {
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

	if r.DurationFromDays != nil {
		v.DurationFromDays = r.DurationFromDays
	}

	if r.DurationToDays != nil {
		v.DurationToDays = r.DurationToDays
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
