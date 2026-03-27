package list

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Request struct {
	Order         Order
	Limit         int
	EncodedCursor string
	Requirements  *Requirements
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

	if r.Order == OrderSalaryAsc || r.Order == OrderSalaryDesc {
		if r.Requirements != nil &&
			r.Requirements.IsPaid != nil &&
			!*r.Requirements.IsPaid {
			return vacancy.ErrInvalidSalaryOrderForUnpaid
		}
	}

	return nil
}

const (
	defaultLimit = 20
	maxLimit     = 100
)
