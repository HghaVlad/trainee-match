package main

import (
	"context"
	"encoding/binary"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/hamba/avro/v2"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/kafka"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

const (
	userCreateSchemaID = 4
	userCreateEvType   = "UserCreated"
)

// produces user created event to kafka,
// used to mock producer to test consumer,
// requires running schema registry, kafka
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	schemaRegCl := schemaregistry.NewClient(cfg.SchemaRegistry)
	schemaLocalReg, err := schemaregistry.NewLocalRegistry(context.Background(), schemaRegCl)
	if err != nil {
		panic(err)
	}

	kprClient, err := kafka.NewClientForProducer(cfg.Kafka)
	if err != nil {
		panic(err)
	}
	kProducer := kafka.NewProducer(cfg.Kafka, kprClient, slog.Default())

	ev := userhr.CreatedEvent{
		EventID:    uuid.New(),
		UserID:     uuid.New(),
		Username:   "JohnPork360",
		Role:       identity.RoleHR,
		Email:      "john.pork360@mail.com",
		OccurredAt: time.Now(),
	}

	schema, err := schemaLocalReg.GetSchemaByID(context.Background(), userCreateSchemaID)
	if err != nil {
		panic(err)
	}

	payload, err := avro.Marshal(schema, ev)
	if err != nil {
		panic(err)
	}

	schemaB := make([]byte, 5)
	binary.BigEndian.PutUint32(schemaB[1:], uint32(userCreateSchemaID))
	//nolint:makezero // test, under control
	payload = append(schemaB, payload...)

	msgs := []outbox.Message{
		{
			ID:        ev.EventID,
			Key:       ev.UserID[:],
			Payload:   payload,
			EventType: userCreateEvType,
			Topic:     cfg.Kafka.UserTopic,
		},
	}

	kProducer.ProduceOutbox(context.Background(), msgs)
}
