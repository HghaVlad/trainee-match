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

func (d *Decoder) DecodeResumeDeletedEvent(_ context.Context, data []byte) (projection.ResumeDeletedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.ResumeDeletedEvent{}, err
	}

	var event projection.ResumeDeletedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.ResumeDeletedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}

	return event, nil
}

func (d *Decoder) DecodeCandidateUpsertedEvent(
	_ context.Context,
	data []byte,
) (projection.CandidateUpsertedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.CandidateUpsertedEvent{}, err
	}

	var event projection.CandidateUpsertedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.CandidateUpsertedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeCompanyUpdatedEvent(_ context.Context, data []byte) (projection.CompanyUpdatedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.CompanyUpdatedEvent{}, err
	}

	var event projection.CompanyUpdatedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.CompanyUpdatedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeCompanyDeletedEvent(_ context.Context, data []byte) (projection.CompanyDeletedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.CompanyDeletedEvent{}, err
	}

	var event projection.CompanyDeletedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.CompanyDeletedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeCompanyMemberAddedEvent(
	_ context.Context,
	data []byte,
) (projection.CompanyMemberAddedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.CompanyMemberAddedEvent{}, err
	}

	var event projection.CompanyMemberAddedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.CompanyMemberAddedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeCompanyMemberRemovedEvent(
	_ context.Context,
	data []byte,
) (projection.CompanyMemberRemovedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.CompanyMemberRemovedEvent{}, err
	}

	var event projection.CompanyMemberRemovedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.CompanyMemberRemovedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeVacancyPublishedEvent(
	_ context.Context,
	data []byte,
) (projection.VacancyPublishedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.VacancyPublishedEvent{}, err
	}

	var event projection.VacancyPublishedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.VacancyPublishedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeVacancyArchivedEvent(
	_ context.Context,
	data []byte,
) (projection.VacancyArchivedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.VacancyArchivedEvent{}, err
	}

	var event projection.VacancyArchivedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.VacancyArchivedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func (d *Decoder) DecodeVacancyUpdatedEvent(
	_ context.Context,
	data []byte,
) (projection.VacancyUpdatedEvent, error) {
	schemaID := getSchemaID(data)
	schema, err := d.registry.GetSchemaByID(schemaID)
	if err != nil {
		return projection.VacancyUpdatedEvent{}, err
	}

	var event projection.VacancyUpdatedEvent
	err = avro.Unmarshal(schema, data[5:], &event)
	if err != nil {
		return projection.VacancyUpdatedEvent{}, fmt.Errorf("failed to unmarshal avro event: %w", err)
	}
	return event, nil
}

func getSchemaID(bytes []byte) int {
	bytes = bytes[1:] // magic byte
	schemaID := binary.BigEndian.Uint32(bytes)
	return int(schemaID)
}
