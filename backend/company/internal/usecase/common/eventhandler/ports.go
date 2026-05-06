package eventhandler

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

//go:generate mockgen -source=ports.go -destination=mocks/mocks.go -package=mocks
type DLQSender interface {
	ToDLQ(ctx context.Context, eventID uuid.UUID, key, payload []byte, topic, eventType string, errMsg string) error
}

type Decoder interface {
	GetUserCreatedEvent(ctx context.Context, payload []byte) (*userhr.CreatedEvent, error)
}

type UserHrCreator interface {
	Execute(ctx context.Context, ev userhr.CreatedEvent) error
}
