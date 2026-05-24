package schemaregistry

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
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

func (d *Decoder) GetUserCreatedEvent(ctx context.Context, payload []byte) (*userhr.CreatedEvent, error) {
	return decodeEvent[userhr.CreatedEvent](ctx, d, payload)
}

func (d *Decoder) GetVacancyPublishedEvent(ctx context.Context, payload []byte) (*vacancy.PublishedEvent, error) {
	return decodeEvent[vacancy.PublishedEvent](ctx, d, payload)
}

func (d *Decoder) GetVacancyDraftCreatedEvent(ctx context.Context, payload []byte) (*vacancy.DraftCreatedEvent, error) {
	return decodeEvent[vacancy.DraftCreatedEvent](ctx, d, payload)
}

func (d *Decoder) GetVacancyUpdatedEvent(ctx context.Context, payload []byte) (*vacancy.UpdatedEvent, error) {
	return decodeEvent[vacancy.UpdatedEvent](ctx, d, payload)
}

func (d *Decoder) GetVacancyArchivedEvent(ctx context.Context, payload []byte) (*vacancy.ArchivedEvent, error) {
	return decodeEvent[vacancy.ArchivedEvent](ctx, d, payload)
}

func (d *Decoder) GetVacancyModUpdEvent(ctx context.Context, payload []byte) (*vacancy.ModerationUpdatedEvent, error) {
	return decodeEvent[vacancy.ModerationUpdatedEvent](ctx, d, payload)
}

func (d *Decoder) GetCompanyUpdatedEvent(ctx context.Context, payload []byte) (*company.UpdatedEvent, error) {
	return decodeEvent[company.UpdatedEvent](ctx, d, payload)
}

func (d *Decoder) GetCompanyDeletedEvent(ctx context.Context, payload []byte) (*company.DeletedEvent, error) {
	return decodeEvent[company.DeletedEvent](ctx, d, payload)
}

func decodeEvent[T any](ctx context.Context, d *Decoder, payload []byte) (*T, error) {
	if len(payload) < magicAndFourBytes {
		return nil, errors.New("missing schema id in avro wire bytes")
	}

	schemaID := getSchemaID(payload)
	payload = payload[magicAndFourBytes:]

	schema, err := d.registry.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return nil, fmt.Errorf("decode event: %w", err)
	}

	var event T

	err = avro.Unmarshal(schema, payload, &event)
	if err != nil {
		return nil, fmt.Errorf("%w: decode event: %w", ErrDecodePayload, err)
	}

	return &event, nil
}

func getSchemaID(bytes []byte) int {
	bytes = bytes[1:] // magic byte
	schemaID := binary.BigEndian.Uint32(bytes)
	return int(schemaID)
}
