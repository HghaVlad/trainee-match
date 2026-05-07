package listhrsummary

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
)

type Request struct {
	Statuses    []application.Status
	CompanyID   *uuid.UUID
	VacancyID   *uuid.UUID
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	Cursor      string
	Limit       int
	Order       cursors.HrSummaryOrder
}

var ErrCompanyOrVacancyRequired = errors.New("either company or vacancy must be set")

func (r *Request) validate() error {
	if r.CompanyID == nil && r.VacancyID == nil {
		return ErrCompanyOrVacancyRequired
	}

	if !r.Order.IsValid() {
		return cursors.ErrUnsupportedOrder
	}

	return nil
}

func (r *Request) normalize() {
	if r.Limit <= 0 || r.Limit > 100 {
		r.Limit = 20
	}

	if r.Order == "" {
		r.Order = cursors.HrSummaryOrderCreatedAtDesc
	}
}
