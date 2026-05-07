package listcandidatesummary

import (
	"time"

	"github.com/google/uuid"
)

type SummaryCursor struct {
	SortAt time.Time `json:"sortAt"`
	AppID  uuid.UUID `json:"appId"`
}

func (o Order) IsValid() bool {
	switch o {
	case OrderCreatedAtDesc, OrderUpdatedAtDesc:
		return true
	default:
		return false
	}
}
