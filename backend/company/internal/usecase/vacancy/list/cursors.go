package list_vacancy

import (
	"time"

	"github.com/google/uuid"
)

type Order string

const (
	OrderPublishedAtDesc Order = "published_at_desc"
	OrderSalaryDesc      Order = "salary_desc"
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
