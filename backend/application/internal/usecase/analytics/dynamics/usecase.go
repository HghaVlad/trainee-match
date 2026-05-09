package dynamics

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	historyRepo historyRepo
	compMemRepo companyMemberRepo
	vacRepo     vacancyProjRepo
}

func NewUsecase(
	historyRepo historyRepo,
	compMemRepo companyMemberRepo,
	vacRepo vacancyProjRepo,
) *Usecase {
	return &Usecase{
		historyRepo: historyRepo,
		compMemRepo: compMemRepo,
		vacRepo:     vacRepo,
	}
}

func (u *Usecase) GetDashboardByCompany(
	ctx context.Context,
	compID uuid.UUID,
	period Period,
	iden identity.Identity,
) ([]Bucket, error) {
	if iden.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	period.Normalize()
	if err := period.Validate(); err != nil {
		return nil, err
	}

	ok, err := u.compMemRepo.IsMember(ctx, iden.UserID, compID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, projection.ErrCompanyNotFound
	}

	dashboard, err := u.historyRepo.GetDynamicsBucketsByCompany(
		ctx, compID,
		*period.From, *period.To, period.Interval,
	)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (u *Usecase) GetDashboardByVacancy(
	ctx context.Context,
	vacID uuid.UUID,
	period Period,
	iden identity.Identity,
) ([]Bucket, error) {
	if iden.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	period.Normalize()
	if err := period.Validate(); err != nil {
		return nil, err
	}

	_, err := u.vacRepo.CheckHrAccess(ctx, iden.UserID, vacID)
	if err != nil {
		return nil, err
	}

	dashboard, err := u.historyRepo.GetDynamicsBucketsByVacancy(
		ctx, vacID,
		*period.From, *period.To, period.Interval,
	)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
