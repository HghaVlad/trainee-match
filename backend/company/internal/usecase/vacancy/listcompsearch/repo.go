package listcompsearch

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type VacancyRepo interface {
	ListByCompanySummaries(
		ctx context.Context,
		requirements *listsearch.Requirements,
		status *vacancy.Status,
		order Order,
		cursor any,
		limit int,
	) (*SearchResult, error)
}

type SearchResult struct {
	Vacancies  []views.MemberVacSummary
	NextCursor any
	HasNext    bool
}

type CompanyRepo interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
