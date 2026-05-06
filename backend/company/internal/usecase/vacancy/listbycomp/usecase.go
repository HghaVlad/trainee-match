package listbycomp

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

	var resp *Response
	var err error

	switch req.Order {
	case OrderCreatedAtDesc:
		resp, err = listByCreatedAt(ctx, uc, req)

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	uc.respCache.Put(ctx, respCacheKey, resp, time.Second*20)
	return resp, nil
}

func listByCreatedAt(ctx context.Context, uc *Usecase, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[CreatedAtCursor, Order](req.EncodedCursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	vacancies, err := uc.vacRepo.ListByCompanySummaries(
		ctx,
		req.CompID,
		req.Requirements,
		req.Status,
		cursor,
		req.Limit+1,
	)
	if err != nil {
		return nil, err
	}

	nextCursor, vacancies := getNextCursor(vacancies, req.Limit)

	resp, err := buildResponse(vacancies, nextCursor, req.Order)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func getNextCursor(vacancies []VacancySummary, limit int) (*CreatedAtCursor, []VacancySummary) {
	if len(vacancies) <= limit {
		return nil, vacancies
	}

	vacancies = vacancies[:len(vacancies)-1]
	last := vacancies[len(vacancies)-1]

	return &CreatedAtCursor{
		CreatedAt: last.CreatedAt,
		ID:        last.ID,
	}, vacancies
}

func buildResponse(vacancies []VacancySummary, nextCursor *CreatedAtCursor, order Order) (*Response, error) {
	nextCursorEncoded, err := encoding.EncodeCursor[CreatedAtCursor, Order](order, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Vacancies:  vacancies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}
