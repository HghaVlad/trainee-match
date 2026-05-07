package cursors

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type SummaryCursor struct {
	SortAt time.Time `json:"sortAt"`
	AppID  uuid.UUID `json:"appId"`
}

type HrSummaryCursor struct {
	SortAt   *time.Time `json:"sortAt,omitempty"`
	FullName *string    `json:"fullName,omitempty"`
	AppID    uuid.UUID  `json:"appId"`
}

func (o SummaryOrder) IsValid() bool {
	switch o {
	case OrderCreatedAtDesc, OrderUpdatedAtDesc:
		return true
	default:
		return false
	}
}

var (
	ErrInvalidCursor       = errors.New("invalid cursor")
	ErrCursorOrderMismatch = errors.New("cursor order mismatch")
	ErrUnsupportedOrder    = errors.New("unsupported list order")
)

type SummaryOrder string

const (
	OrderCreatedAtDesc = "createdAtDesc"
	OrderUpdatedAtDesc = "updatedAtDesc"
)

type HrSummaryOrder string

const (
	HrSummaryOrderCreatedAtDesc     = "createdAtDesc"
	HrSummaryOrderUpdatedAtDesc     = "updatedAtDesc"
	HrSummaryOrderCandidateFullName = "candidateFullNameAsc"
)

func (o HrSummaryOrder) IsValid() bool {
	switch o {
	case HrSummaryOrderCreatedAtDesc, HrSummaryOrderUpdatedAtDesc, HrSummaryOrderCandidateFullName:
		return true
	default:
		return false
	}
}
