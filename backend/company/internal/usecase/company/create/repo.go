package create_company

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type CompanyRepo interface {
	Create(ctx context.Context, company *domain.Company) error
}

type CompanyMemberRepo interface {
	Create(ctx context.Context, member *domain.CompanyMember) error
}
