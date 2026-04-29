package listbycomp

import (
	"time"

	"github.com/google/uuid"
)

type Order string

const (
	OrderCreatedAtDesc Order = "created_at_desc"
)

type CreatedAtCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

func (r Order) IsValid() bool {
	return r == OrderCreatedAtDesc
}
