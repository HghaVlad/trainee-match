package listbycomp

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	vaclist "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type VacancyRepo interface {
	ListByCompanySummaries(
		ctx context.Context,
		compID uuid.UUID,
		requirements *vaclist.Requirements,
		status *vacancy.Status,
		cursor *CreatedAtCursor,
		limit int,
	) ([]VacancySummary, error)
}

type CompanyRepo interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
