package listcompsearch

import (
	"time"

	"github.com/google/uuid"
)

type Order string

const (
	OrderCreatedAtDesc Order = "created_at_desc"
	OrderRelevance     Order = "relevance"
)

type RelevanceCursor struct {
	Relevance float64
	CreatedAt time.Time
	ID        uuid.UUID
}

type CreatedAtCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

func (r Order) IsValid() bool {
	switch r {
	case OrderRelevance, OrderCreatedAtDesc:
		return true
	default:
		return false
	}
}
