package outbox

import "context"

type WriterRepo interface {
	Create(ctx context.Context, msg Message) error
}
