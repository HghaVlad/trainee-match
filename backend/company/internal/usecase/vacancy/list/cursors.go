package list_vacancy

import (
	"time"

	"github.com/google/uuid"
)

type Order string

const (
	OrderPublishedAtDesc Order = "published_at_desc"
	OrderSalaryDesc      Order = "salary_desc"
	OrderSalaryAsc       Order = "salary_asc"
)

type PublishedAtCursor struct {
	PublishedAt time.Time
	Id          uuid.UUID
}

type SalaryCursor struct {
	SalaryFrom int
	SalaryTo   int
	Id         uuid.UUID
}

func (r Order) IsValid() bool {
	switch r {
	case OrderPublishedAtDesc,
		OrderSalaryDesc,
		OrderSalaryAsc:
		return true
	}

	return false
}
