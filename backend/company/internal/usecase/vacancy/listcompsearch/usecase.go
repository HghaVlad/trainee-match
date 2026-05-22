package listcompsearch

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/encoding"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type Usecase struct {
	vacRepo    VacancyRepo
	compRepo   CompanyRepo
	memberRepo CompMemberRepo
	respCache  ResponseCacheRepo
}

func NewUsecase(
	vacRepo VacancyRepo,
	compRepo CompanyRepo,
	memberRepo CompMemberRepo,
	cache ResponseCacheRepo,
) *Usecase {
	return &Usecase{
		vacRepo:    vacRepo,
		memberRepo: memberRepo,
		compRepo:   compRepo,
		respCache:  cache,
	}
}

func (uc *Usecase) Execute(ctx context.Context, req *Request, ident *identity.Identity) (*Response, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := uc.authorize(ctx, req.CompID, ident); err != nil {
		return nil, err
	}

	respCacheKey := requestToCacheKey(req)
	res := uc.respCache.Get(ctx, respCacheKey)
	if res != nil {
		return res, nil
	}

	companyExists, exErr := uc.compRepo.Exists(ctx, req.CompID)
	if exErr != nil {
		return nil, exErr
	}

	if !companyExists {
		return nil, company.ErrCompanyNotFound
	}

	req.Requirements.Companies = &[]uuid.UUID{req.CompID}

	if (req.Requirements.Query == nil || *req.Requirements.Query == "") &&
		req.Order == OrderRelevance {
		req.Order = OrderCreatedAtDesc
		req.EncodedCursor = ""
	}

	var resp *Response
	var err error

	switch req.Order {
	case OrderRelevance:
		resp, err = list[RelevanceCursor](ctx, uc, req)
	case OrderCreatedAtDesc:
		resp, err = list[CreatedAtCursor](ctx, uc, req)

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	uc.respCache.Put(ctx, respCacheKey, resp, time.Second*20)
	return resp, nil
}

func list[CursorT any](ctx context.Context, uc *Usecase, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[CursorT, Order](req.EncodedCursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	res, err := uc.vacRepo.ListByCompanySummaries(ctx, req.Requirements, req.Status, req.Order, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	if res.HasNext {
		nextCursor, _ := res.NextCursor.(*CursorT)
		return buildResponse[CursorT](res.Vacancies, nextCursor, req.Order)
	}

	return buildResponse[CursorT](res.Vacancies, nil, req.Order)
}

// only member of company can view their vacancies in full
func (uc *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident *identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	_, err := uc.memberRepo.Get(ctx, ident.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}

	return err
}

func buildResponse[CursorT any](
	vacancies []views.MemberVacSummary,
	nextCursor *CursorT,
	order Order,
) (*Response, error) {
	nextCursorEncoded, err := encoding.EncodeCursor[CursorT, Order](order, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Vacancies:  vacancies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}
