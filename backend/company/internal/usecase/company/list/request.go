package list

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Request struct {
	Order         Order
	Limit         int
	EncodedCursor string
	Filter        Filter
}

type Filter struct {
	CompanyMemberID *uuid.UUID
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

	return nil
}

const (
	defaultLimit = 20
	maxLimit     = 100
)

func (r Order) IsValid() bool {
	switch r {
	case OrderVacanciesDesc,
		OrderNameAsc,
		OrderCreatedAtDesc:
		return true
	}

	return false
}
