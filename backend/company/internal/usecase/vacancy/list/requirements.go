package list_vacancy

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/google/uuid"
)

type RangeInt struct {
	Min *int
	Max *int
}

type Requirements struct {
	Salary            *RangeInt
	HoursPerWeek      *RangeInt
	Duration          *RangeInt
	WorkFormat        *[]value_types.WorkFormat
	Companies         *[]uuid.UUID
	City              *[]string
	IsPaid            *bool
	InternshipToOffer *bool
	FlexibleSchedule  *bool
}
