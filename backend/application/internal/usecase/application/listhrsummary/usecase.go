package listhrsummary

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
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
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
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
		ctx, ident.UserID, req.Statuses,
		req.CompanyID, req.VacancyID,
		req.CreatedFrom, req.CreatedTo,
		req.Order, cursor, req.Limit+1,
	)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		if err := u.authorize(ctx, req.CompanyID, req.VacancyID, ident); err != nil {
			return nil, err
		}
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
	if vacID != nil {
		realCompID, err := u.vacRepo.CheckHrAccess(ctx, ident.UserID, *vacID)
		if err != nil {
			return err
		}

		if compID != nil && *compID != realCompID {
			return projection.ErrVacancyNotFound
		}

		return nil
	}

	if compID != nil {
		ok, err := u.memRepo.IsMember(ctx, ident.UserID, *compID)
		if err != nil {
			return err
		}

		if !ok {
			return projection.ErrCompanyNotFound
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
