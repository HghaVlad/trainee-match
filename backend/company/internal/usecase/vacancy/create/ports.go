package create

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

//go:generate mockgen -source=ports.go -destination=mocks/port_mocks.go -package=mocks
type VacancyRepo interface {
	Create(ctx context.Context, vacancy *vacancy.Vacancy) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error)
}

type CompanyRepo interface {
	GetByMember(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*company.Company, error)
}

type outboxWriter interface {
	WriteVacancyDraftCreated(ctx context.Context, ev vacancy.DraftCreatedEvent) error
}
