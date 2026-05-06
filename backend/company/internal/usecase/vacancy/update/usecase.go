package update

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	repo       VacancyRepo
	memberRepo CompMemberRepo
	outbox     outboxWriter
	cache      CacheRepo
	txManager  common.TxManager
}

func NewUsecase(
	repo VacancyRepo,
	memberRepo CompMemberRepo,
	outbox outboxWriter,
	cacheRepo CacheRepo,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		repo:       repo,
		memberRepo: memberRepo,
		outbox:     outbox,
		cache:      cacheRepo,
		txManager:  txManager,
	}
}

// Execute updates vacancy. All nil fields of vacancy in request won't be applied.
// Deletes vacancy from cache.
func (u *Usecase) Execute(ctx context.Context, req *Request, identity *identity.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if err := req.lightValidate(); err != nil {
		return err
	}

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, req.CompanyID, identity); err != nil {
			return err
		}

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

		if eventShouldBeCreated {
			return u.createdVacancyUpdatedEvent(ctx, vac)
		}

		return nil
	})
	if err != nil {
		return err
	}

	u.cache.Del(ctx, req.VacancyID)
	return nil
}

// only member of company can update vacancy
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident *identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, ident.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}

	return err
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
