package listcandidatesummary

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	repo repo
}

func NewUsecase(repo repo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(ctx context.Context, req Request, ident identity.Identity) (*Response, error) {
	if ident.Role != identity.RoleCandidate {
		return nil, identity.ErrCandidateRoleRequired
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	if req.Order == "" {
		req.Order = cursors.OrderCreatedAtDesc
	}

	if !req.Order.IsValid() {
		return nil, cursors.ErrUnsupportedOrder
	}

	cursor, err := cursors.DecodeSummaryCursor(req.Cursor, req.Order)
	if err != nil {
		return nil, err
	}

	items, err := u.repo.ListCandidateAppSummaries(
		ctx,
		ident.UserID,
		req.Statuses,
		req.CompanyID,
		req.Order,
		cursor,
		req.Limit+1,
	)
	if err != nil {
		return nil, err
	}

	nextCursor, items := getNextCursor(items, req.Limit, req.Order)
	encodedCursor, err := cursors.EncodeCursor(req.Order, nextCursor)
	if err != nil {
		return nil, err
	}

	return &Response{
		AppSummaries: items,
		NextCursor:   encodedCursor,
		HasNext:      nextCursor != nil,
	}, nil
}

func getNextCursor(
	items []views.CandidateAppSummary,
	limit int,
	order cursors.SummaryOrder,
) (*cursors.SummaryCursor, []views.CandidateAppSummary) {
	if len(items) <= limit {
		return nil, items
	}

	items = items[:limit]
	last := items[len(items)-1]

	var sortAt time.Time
	switch order {
	case cursors.OrderUpdatedAtDesc:
		sortAt = last.UpdatedAt
	default:
		sortAt = last.CreatedAt
	}

	return &cursors.SummaryCursor{
		SortAt: sortAt,
		AppID:  last.AppID,
	}, items
}
