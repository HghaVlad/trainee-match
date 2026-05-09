package summary

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	appRepo           appRepo
	companyMemberRepo companyMemberRepo
	vacRepo           vacancyProjRepo
}

func NewUsecase(repo appRepo, companyMemberRepo companyMemberRepo, vacRepo vacancyProjRepo) *Usecase {
	return &Usecase{
		appRepo:           repo,
		companyMemberRepo: companyMemberRepo,
		vacRepo:           vacRepo,
	}
}

func (u *Usecase) GetByCompany(
	ctx context.Context,
	compID uuid.UUID,
	ident identity.Identity,
) (*Summary, error) {
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	ok, err := u.companyMemberRepo.IsMember(ctx, ident.UserID, compID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, projection.ErrCompanyNotFound
	}

	sum, err := u.appRepo.GetCompanyAnalyticsSummary(ctx, compID)
	if err != nil {
		return nil, err
	}

	enrichSummary(sum)
	return sum, nil
}

func (u *Usecase) GetByVacancy(
	ctx context.Context,
	vacID uuid.UUID,
	ident identity.Identity,
) (*Summary, error) {
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	_, err := u.vacRepo.CheckHrAccess(ctx, ident.UserID, vacID)
	if err != nil {
		return nil, err
	}

	sum, err := u.appRepo.GetVacancyAnalyticsSummary(ctx, vacID)
	if err != nil {
		return nil, err
	}

	enrichSummary(sum)
	return sum, nil
}

func enrichSummary(sum *Summary) {
	sum.ActiveApplications = sum.SubmittedCount + sum.SeenCount + sum.InterviewCount

	sum.TotalApplications = sum.ActiveApplications + sum.OfferCount +
		sum.RejectedCount + sum.WithdrawnCount

	if sum.TotalApplications > 0 {
		sum.ConversionToInterview =
			float32(sum.InterviewCount) / float32(sum.TotalApplications)

		sum.ConversionToOffer =
			float32(sum.OfferCount) / float32(sum.TotalApplications)
	}
}
