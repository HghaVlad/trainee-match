package schemaregistry

import (
	"encoding/binary"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

type Encoder struct {
	registry *LocalRegistry
}

func NewEncoder(registry *LocalRegistry) *Encoder {
	return &Encoder{registry}
}

func (en *Encoder) UserCreatedToBytes(event domain.UserCreatedEvent) ([]byte, error) {
	return en.eventToBytes(event, "user-created-value")
}

func (en *Encoder) eventToBytes(ev any, subject string) ([]byte, error) {
	schemaID, err := en.registry.GetSchemaIDBySubject(subject)
	if err != nil {
		return nil, err
	}
	schema, err := en.registry.GetSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	payload, err := avro.Marshal(schema, ev)
	if err != nil {
		return nil, fmt.Errorf("event to bytes for subject %s: avro marshal: %w", subject, err)
	}

	const magicByteAndUint32 = 5
	bytes := make([]byte, magicByteAndUint32+len(payload))

	writeConfluentWireSchemaID(bytes, schemaID)

	copy(bytes[magicByteAndUint32:], payload)
	return bytes, nil
}

// writes [0][schemaID] - payload to be appended
func writeConfluentWireSchemaID(buf []byte, schemaID int) {
	buf[0] = 0 // magic byte
	//nolint:gosec // under control
	binary.BigEndian.PutUint32(buf[1:], uint32(schemaID))
}
