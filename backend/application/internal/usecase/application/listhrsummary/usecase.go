package listhrsummary

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	appRepo appRepo
	memRepo memProjRepo
	vacRepo vacProjRepo
}

func NewUsecase(repo appRepo, memRepo memProjRepo, vacRepo vacProjRepo) *Usecase {
	return &Usecase{
		appRepo: repo,
		memRepo: memRepo,
		vacRepo: vacRepo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req Request, ident identity.Identity) (*Response, error) {
	if err := u.authorize(ctx, req.CompanyID, req.VacancyID, ident); err != nil {
		return nil, err
	}

	req.normalize()
	if err := req.validate(); err != nil {
		return nil, err
	}

	cursor, err := cursors.DecodeHrSummaryCursor(req.Cursor, req.Order)
	if err != nil {
		return nil, err
	}

	items, err := u.appRepo.ListHrAppSummaries(
		ctx,
		req.Statuses,
		req.CompanyID,
		req.VacancyID,
		req.CreatedFrom,
		req.CreatedTo,
		req.Order,
		cursor,
		req.Limit+1,
	)
	if err != nil {
		return nil, err
	}

	nextCursor, items := getNextCursor(items, req.Limit, req.Order)
	encodedCursor, err := cursors.EncodeHrSummaryCursor(req.Order, nextCursor)
	if err != nil {
		return nil, err
	}

	return &Response{
		AppSummaries: items,
		NextCursor:   encodedCursor,
		HasNext:      nextCursor != nil,
	}, nil
}

func (u *Usecase) authorize(ctx context.Context, compID, vacID *uuid.UUID, ident identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	if vacID != nil {
		cID, err := u.vacRepo.GetCompanyIDByVacancyID(ctx, *vacID)
		if err != nil {
			return err
		}

		compID = &cID
	}

	if compID != nil {
		ok, err := u.memRepo.IsMember(ctx, ident.UserID, *compID)
		if err != nil {
			return err
		}

		if !ok {
			return application.ErrAccessDenied
		}

		return nil
	}

	return nil
}

func getNextCursor(
	items []views.HrAppSummary,
	limit int,
	order cursors.HrSummaryOrder,
) (*cursors.HrSummaryCursor, []views.HrAppSummary) {
	if len(items) <= limit {
		return nil, items
	}

	items = items[:limit]
	last := items[len(items)-1]

	switch order {
	case cursors.HrSummaryOrderUpdatedAtDesc:
		sortAt := last.UpdatedAt
		return &cursors.HrSummaryCursor{
			SortAt: &sortAt,
			AppID:  last.AppID,
		}, items
	case cursors.HrSummaryOrderCandidateFullName:
		fullName := last.AppSnap.FullName
		return &cursors.HrSummaryCursor{
			FullName: &fullName,
			AppID:    last.AppID,
		}, items
	default:
		sortAt := last.CreatedAt
		return &cursors.HrSummaryCursor{
			SortAt: &sortAt,
			AppID:  last.AppID,
		}, items
	}
}
