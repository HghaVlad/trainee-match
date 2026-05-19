package create

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type CompanyRepo interface {
	Create(ctx context.Context, company *company.Company) error
}

type CompanyMemberRepo interface {
	Create(ctx context.Context, member *member.CompanyMember) error
}

type outboxWriter interface {
	WriteCompanyMemberAdded(ctx context.Context, ev member.AddedEvent) error
}
