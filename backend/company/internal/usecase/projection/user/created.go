package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type CreatedEvent struct {
	EventID    uuid.UUID           `avro:"event_id"`
	UserID     uuid.UUID           `avro:"user_id"`
	Username   string              `avro:"username"`
	Role       identity.GlobalRole `avro:"role"`
	Email      string              `avro:"email"`
	OccurredAt time.Time           `avro:"occurred_at"`
}
