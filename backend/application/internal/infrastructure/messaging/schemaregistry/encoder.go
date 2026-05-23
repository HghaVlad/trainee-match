package schemaregistry

import (
	"encoding/binary"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/dlq"
)

const (
	dlqSubject = "dlq-value"
)

type Encoder struct {
	registry *LocalRegistry
}

func NewEncoder(registry *LocalRegistry) *Encoder {
	return &Encoder{registry: registry}
}

func (en *Encoder) DLQToBytes(message dlq.Message) ([]byte, error) {
	return en.eventToBytes(dlqSubject, message)
}

func (en *Encoder) eventToBytes(subject string, event any) ([]byte, error) {
	schemaID, err := en.registry.GetSchemaIDBySubject(subject)
	if err != nil {
		return nil, err
	}
	schema, err := en.registry.GetSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	encodedEvent, err := avro.Marshal(schema, event)
	if err != nil {
		return nil, fmt.Errorf("error encoding event for subject '%s': %w", subject, err)
	}

	const magicByteAndUint32 = 5
	bytes := make([]byte, len(encodedEvent)+magicByteAndUint32)

	writeConfluentWireSchemaID(bytes, schemaID)
	copy(bytes[magicByteAndUint32:], encodedEvent)

	return bytes, nil
}

// writes [0][schemaID] - payload to be appended
func writeConfluentWireSchemaID(buf []byte, schemaID int) {
	buf[0] = 0 // magic byte
	//nolint:gosec // under control
	binary.BigEndian.PutUint32(buf[1:], uint32(schemaID))
}
