package schemaregistry

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Decoder struct {
	registry *LocalRegistry
}

func NewDecoder(registry *LocalRegistry) *Decoder {
	return &Decoder{registry: registry}
}

func (d *Decoder) DecodeResumeUpsertedEvent(_ context.Context, data []byte) (projection.ResumeUpsertedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.ResumeUpsertedEvent{}, err
	}

	var event projection.ResumeUpsertedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.ResumeUpsertedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}

	return event, nil
}

func getSchemaID(bytes []byte) int {
	bytes = bytes[1:] // magic byte
	schemaID := binary.BigEndian.Uint32(bytes)
	return int(schemaID)
}
