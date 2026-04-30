package schemaregistry

import (
	"encoding/binary"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

const (
	vacancyPublishedSubject = "vacancy-published-value"
	vacancyArchivedSubject  = "vacancy-archived-value"
	vacancyUpdatedSubject   = "vacancy-updated-value"

	companyUpdatedSubject = "company-updated-value"
	companyDeletedSubject = "company-deleted-value"

	recruiterAddedSubject   = "company-member-added-value"
	recruiterRemovedSubject = "company-member-removed-value"
)

type Encoder struct {
	localReg *LocalRegistry
	codecs   map[string]avro.Schema
}

func NewEncoder(localReg *LocalRegistry) (*Encoder, error) {
	codecs := make(map[string]avro.Schema)

	for subj, rawSchema := range localReg.schemas {
		schema, err := avro.Parse(rawSchema)
		if err != nil {
			return nil, fmt.Errorf("new schema reg encoder: avro parse schema: %w", err)
		}

		codecs[subj] = schema
	}

	return &Encoder{
		localReg: localReg,
		codecs:   codecs,
	}, nil
}

func (en *Encoder) VacancyPublishedToBytes(ev vacancy.PublishedEvent) ([]byte, error) {
	return en.eventToBytes(ev, vacancyPublishedSubject)
}

func (en *Encoder) VacancyArchivedToBytes(ev vacancy.ArchivedEvent) ([]byte, error) {
	return en.eventToBytes(ev, vacancyArchivedSubject)
}

func (en *Encoder) VacancyUpdatedToBytes(ev vacancy.UpdatedEvent) ([]byte, error) {
	return en.eventToBytes(ev, vacancyUpdatedSubject)
}

func (en *Encoder) CompanyMemberAddedToBytes(ev member.RecruiterAddedEvent) ([]byte, error) {
	return en.eventToBytes(ev, recruiterAddedSubject)
}

func (en *Encoder) CompanyMemberRemovedToBytes(ev member.RecruiterRemovedEvent) ([]byte, error) {
	return en.eventToBytes(ev, recruiterRemovedSubject)
}

func (en *Encoder) CompanyDeletedToBytes(ev company.DeletedEvent) ([]byte, error) {
	return en.eventToBytes(ev, companyDeletedSubject)
}

func (en *Encoder) CompanyUpdatedToBytes(ev company.UpdatedEvent) ([]byte, error) {
	return en.eventToBytes(ev, companyUpdatedSubject)
}

func (en *Encoder) eventToBytes(ev any, subject string) ([]byte, error) {
	codec, ok := en.codecs[subject]
	if !ok {
		return nil, fmt.Errorf("event to bytes: missing schema for subject %s", subject)
	}

	schemaID, ok := en.localReg.schemaIDs[subject]
	if !ok {
		return nil, fmt.Errorf("event to bytes: missing schemaID for subject %s", subject)
	}

	payload, err := avro.Marshal(codec, ev)
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
