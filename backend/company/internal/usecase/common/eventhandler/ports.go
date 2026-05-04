package eventhandler

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/user"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
)

//go:generate mockgen -source=handler_ports.go -destination=mocks/mocks.go -package=mocks
type outboxDLQWriter interface {
	WriteToDLQ(ctx context.Context, meta outbox.DLQMeta) error
}

type decoder interface {
	GetUserCreatedEvent(ctx context.Context, schemaID int, allBytes []byte) (*user.CreatedEvent, error)
	RetrieveSchemaID(bytes []byte) (int, error) // TODO: think if we need it
}
