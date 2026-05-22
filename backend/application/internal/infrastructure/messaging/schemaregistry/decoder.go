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

func (d *Decoder) DecodeResumeUpsertedEvent(ctx context.Context, data []byte) (projection.ResumeUpsertedEvent, error) {
	var event projection.ResumeUpsertedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeResumeDeletedEvent(ctx context.Context, data []byte) (projection.ResumeDeletedEvent, error) {
	var event projection.ResumeDeletedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeCandidateUpsertedEvent(
	ctx context.Context,
	data []byte,
) (projection.CandidateUpsertedEvent, error) {
	var event projection.CandidateUpsertedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeCompanyUpdatedEvent(ctx context.Context, data []byte) (projection.CompanyUpdatedEvent, error) {
	var event projection.CompanyUpdatedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeCompanyDeletedEvent(ctx context.Context, data []byte) (projection.CompanyDeletedEvent, error) {
	var event projection.CompanyDeletedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeCompanyMemberAddedEvent(
	ctx context.Context,
	data []byte,
) (projection.CompanyMemberAddedEvent, error) {
	var event projection.CompanyMemberAddedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeCompanyMemberRemovedEvent(
	ctx context.Context,
	data []byte,
) (projection.CompanyMemberRemovedEvent, error) {
	var event projection.CompanyMemberRemovedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeVacancyPublishedEvent(
	ctx context.Context,
	data []byte,
) (projection.VacancyPublishedEvent, error) {
	var event projection.VacancyPublishedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeVacancyArchivedEvent(
	ctx context.Context,
	data []byte,
) (projection.VacancyArchivedEvent, error) {
	var event projection.VacancyArchivedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) DecodeVacancyUpdatedEvent(
	ctx context.Context,
	data []byte,
) (projection.VacancyUpdatedEvent, error) {
	var event projection.VacancyUpdatedEvent
	err := d.decodeEvent(ctx, data, &event)
	return event, err
}

func (d *Decoder) decodeEvent(ctx context.Context, data []byte, event any) error {
	if len(data) < 5 {
		return fmt.Errorf("data is too short data: %v", data)
	}

	schemaID := getSchemaID(data)
	schema, err := d.registry.GetRemoteSchemaByID(ctx, schemaID)
	if err != nil {
		return err
	}

	err = avro.Unmarshal(schema, data[5:], event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal avro event: %w", err)
	}

	return nil
}

func getSchemaID(bytes []byte) int {
	bytes = bytes[1:] // magic byte
	schemaID := binary.BigEndian.Uint32(bytes)
	return int(schemaID)
}
