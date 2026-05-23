package kafka

import (
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
)

func NewClientForProducer(cfg config.Kafka) (*kgo.Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
		kgo.ProducerLinger(cfg.ProducerLinger),
		kgo.RequiredAcks(parseAcks(cfg.ProducerAcks)),
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create franz-go client for producer: %w", err)
	}
	return client, nil
}

func NewClientForConsumer(cfg config.Kafka, consumer *Consumer) (*kgo.Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),

		kgo.ConsumerGroup(cfg.ConsumerGroup),
		kgo.ConsumeTopics(cfg.ConsumerTopics...),
		kgo.DisableAutoCommit(),
		kgo.OnPartitionsAssigned(consumer.Assigned),
		kgo.OnPartitionsRevoked(consumer.Revoked),
		kgo.OnPartitionsLost(consumer.Revoked),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create franz-go client for consumer: %w", err)
	}
	return client, nil
}

func parseAcks(acks string) kgo.Acks {
	switch strings.TrimSpace(strings.ToLower(acks)) {
	case "none", "0":
		return kgo.NoAck()
	case "leader", "1":
		return kgo.LeaderAck()
	case "all", "-1":
		return kgo.AllISRAcks()
	default:
		return kgo.AllISRAcks()
	}
}
