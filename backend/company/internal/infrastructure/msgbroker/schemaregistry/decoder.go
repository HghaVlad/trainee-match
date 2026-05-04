package schemaregistry

import (
	"context"
	"errors"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/user"
)

type Decoder struct {
	registry *LocalRegistry
}

func NewDecoder(registry *LocalRegistry) *Decoder {
	return &Decoder{registry: registry}
}

func (d *Decoder) GetUserCreatedEnvelope(
	ctx context.Context,
	schemaID int,
	allBytes []byte,
) (*UserCreatedEnvelope, error) {
	if len(allBytes) < 5 {
		return nil, errors.New("missing schema id in avro wire bytes")
	}
	payload := allBytes[5:]

	schema, err := d.registry.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return nil, fmt.Errorf("decoder get user created envelope: %w", err)
	}

	var event user.CreatedEvent

	err = avro.Unmarshal(schema, payload, &event)
	if err != nil {
		return nil, fmt.Errorf("decoder get user created envelope: avro unmarshal: %w", err)
	}

	return &UserCreatedEnvelope{
		EventID:  &event.EventID,
		SchemaID: &schemaID,      // TODO: think, if there is a point here
		Event:    &event,
	}, nil
}
