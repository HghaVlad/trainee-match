package list_vacancy

import "github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"

type Request struct {
	Order         Order
	Limit         int
	EncodedCursor string
	Requirements  *Requirements
}

func (r *Request) Validate() error {
	if !r.Order.IsValid() {
		return domain_errors.ErrUnsupportedListOrder
	}

	if r.Limit <= 0 {
		r.Limit = defaultLimit
	}
	if r.Limit > maxLimit {
		return domain_errors.ErrLimitTooLarge
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
			return domain_errors.ErrInvalidSalaryOrderForUnpaid
		}
	}

	return nil
}

const (
	defaultLimit = 20
	maxLimit     = 100
)
