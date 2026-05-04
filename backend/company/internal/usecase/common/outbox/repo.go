package outbox

import "context"

type WriterRepo interface {
	Create(ctx context.Context, msg Message) error
	CreateFailed(ctx context.Context, msg Message) error
}

type RelayRepo interface {
	ListPending(ctx context.Context, limit int) ([]Message, error)
	Save(ctx context.Context, msgs []Message) error
}
