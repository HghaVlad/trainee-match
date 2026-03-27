package list

import (
	"time"
)

type Order string

const (
	OrderCreatedAtDesc Order = "created_at_desc"
	OrderNameAsc       Order = "name_asc"
	OrderVacanciesDesc Order = "vacancies_desc"
)

type CreatedAtCursor struct {
	CreatedAt time.Time
	Name      string
}

type NameCursor struct {
	Name string
}

type VacanciesCntCursor struct {
	Count int
	Name  string
}
