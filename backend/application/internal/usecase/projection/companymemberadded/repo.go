package companymemberadded

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type CompanyMemberRepo interface {
	Save(ctx context.Context, member *projection.CompanyMember) error
}
