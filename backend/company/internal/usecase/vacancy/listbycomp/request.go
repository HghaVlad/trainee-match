package listbycomp

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	vaclist "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type Request struct {
	CompID        uuid.UUID
	Order         Order
	Limit         int
	EncodedCursor string
	Requirements  *vaclist.Requirements
	Status        *vacancy.Status
}

func (r *Request) Validate() error {
	if !r.Order.IsValid() {
		return common.ErrUnsupportedListOrder
	}

	if r.Limit <= 0 {
		r.Limit = defaultLimit
	}
	if r.Limit > maxLimit {
		return common.ErrLimitTooLarge
	}

	if r.Requirements != nil {
		if err := r.Requirements.Validate(); err != nil {
			return err
		}
	}

	if r.Status != nil && !r.Status.IsValid() {
		return vacancy.ErrInvalidStatus
	}

	return nil
}

const (
	defaultLimit = 20
	maxLimit     = 100
)
