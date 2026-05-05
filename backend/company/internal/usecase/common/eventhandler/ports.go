package eventhandler

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/user"
)

//go:generate mockgen -source=handler_ports.go -destination=mocks/mocks.go -package=mocks
type DLQSender interface {
	ToDLQ(ctx context.Context, eventID uuid.UUID, key, payload []byte, topic, eventType string, errMsg string) error
}

type Decoder interface {
	GetUserCreatedEvent(ctx context.Context, payload []byte) (*user.CreatedEvent, error)
}
