package schemaregistry

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/user"
)

type UserCreatedEnvelope struct {
	EventID  *uuid.UUID
	SchemaID *int
	Event    *user.CreatedEvent
}
