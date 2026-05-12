package schemaregistry

import (
	"encoding/binary"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"
)

var (
	ResumeUpsertedEvent    = "resume-upserted-value"
	CandidateUpsertedEvent = "candidate-upserted-value"
	ResumeDeletedEvent     = "resume-deleted-value"
	CandidateDeletedEvent  = "candidate-deleted-value"
)

type Encoder struct {
	registry *LocalRegistry
}

func NewEncoder(registry *LocalRegistry) *Encoder {
	return &Encoder{registry}
}

func (en *Encoder) ResumeUpdatedToBytes(ev events.ResumeUpserted) ([]byte, int, error) {
	return en.EventToBytes(ResumeUpsertedEvent, ev)
}

func (en *Encoder) CandidateUpsertedToBytes(ev events.CandidateUpserted) ([]byte, int, error) {
	return en.EventToBytes(CandidateUpsertedEvent, ev)
}

func (en *Encoder) ResumeDeletedToBytes(ev events.ResumeDeleted) ([]byte, int, error) {
	return en.EventToBytes(ResumeDeletedEvent, ev)
}

func (en *Encoder) CandidateDeletedToBytes(ev events.CandidateDeleted) ([]byte, int, error) {
	return en.EventToBytes(CandidateDeletedEvent, ev)
}

func (en *Encoder) EventToBytes(subject string, event any) ([]byte, int, error) {
	schemaID, err := en.registry.GetSchemaIDBySubject(subject)
	if err != nil {
		return nil, 0, err
	}
	schema, err := en.registry.GetSchemaByID(schemaID)
	if err != nil {
		return nil, 0, err
	}

	encodedEvent, err := avro.Marshal(schema, event)
	if err != nil {
		return nil, 0, err
	}

	const magicByteAndUint32 = 5
	bytes := make([]byte, magicByteAndUint32+len(encodedEvent))

	writeConfluentWireSchemaID(bytes, schemaID)
	copy(bytes[magicByteAndUint32:], encodedEvent)

	return bytes, schemaID, nil
}

// writes [0][schemaID] - payload to be appended
func writeConfluentWireSchemaID(buf []byte, schemaID int) {
	buf[0] = 0 // magic byte
	//nolint:gosec // under control
	binary.BigEndian.PutUint32(buf[1:], uint32(schemaID))
}
