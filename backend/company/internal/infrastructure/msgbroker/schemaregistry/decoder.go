package schemaregistry

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/user"
)

var (
	ErrDecodePayload = errors.New("decode payload error")
)

type Decoder struct {
	registry *LocalRegistry
}

func NewDecoder(registry *LocalRegistry) *Decoder {
	return &Decoder{registry: registry}
}

func (d *Decoder) GetUserCreatedEvent(ctx context.Context, schemaID int, allBytes []byte) (*user.CreatedEvent, error) {
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
		return nil, fmt.Errorf("%v: get user created envelope: %w", ErrDecodePayload, err)
	}

	return &event, nil
}

func (d *Decoder) RetrieveSchemaID(bytes []byte) (int, error) {
	if len(bytes) < 5 {
		return -1, errors.New("missing schema id in avro wire bytes")
	}

	bytes = bytes[1:] // magic byte
	schemaID := binary.BigEndian.Uint32(bytes)
	return int(schemaID), nil
}
