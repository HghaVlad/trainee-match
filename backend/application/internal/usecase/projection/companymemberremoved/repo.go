package companymemberremoved

import (
	"context"

	"github.com/google/uuid"
)

type CompanyMemberRepo interface {
	Delete(ctx context.Context, userID, companyID uuid.UUID) error
}
