package schemaregistry

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/user"
)

const magicAndFourBytes = 5

var (
	ErrDecodePayload = errors.New("decode payload error")
)

type Decoder struct {
	registry *LocalRegistry
}

func NewDecoder(registry *LocalRegistry) *Decoder {
	return &Decoder{registry: registry}
}

func (d *Decoder) GetUserCreatedEvent(ctx context.Context, payload []byte) (*user.CreatedEvent, error) {
	if len(payload) < magicAndFourBytes {
		return nil, errors.New("missing schema id in avro wire bytes")
	}

	schemaID := getSchemaID(payload)
	payload = payload[magicAndFourBytes:]

	schema, err := d.registry.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return nil, fmt.Errorf("decode user created: %w", err)
	}

	var event user.CreatedEvent

	err = avro.Unmarshal(schema, payload, &event)
	if err != nil {
		return nil, fmt.Errorf("%w: decode user created: %w", ErrDecodePayload, err)
	}

	return &event, nil
}

func getSchemaID(bytes []byte) int {
	bytes = bytes[1:] // magic byte
	schemaID := binary.BigEndian.Uint32(bytes)
	return int(schemaID)
}
