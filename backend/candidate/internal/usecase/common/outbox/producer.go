package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Producer interface {
	ProduceOutBox(ctx context.Context, messages []Message) []ProduceResult
}

type ProduceResult struct {
	MsgID       uuid.UUID
	SentAt      *time.Time
	Err         error
	Unretryable bool
}
