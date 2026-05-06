package list

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type Request struct {
	CompanyID uuid.UUID
	Limit     int
	Offset    int
}

func (r *Request) IsValid() error {
	if r.Limit <= 0 {
		r.Limit = defaultLimit
	}
	if r.Limit > maxLimit {
		return common.ErrLimitTooLarge
	}

	return nil
}
