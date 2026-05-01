package kafka

import (
	"errors"
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

func NewClientForProducer(cfg config.Kafka) (*kgo.Client, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),

		kgo.RequiredAcks(parseAcks(cfg.ProducerAcks)),
		kgo.ProducerLinger(cfg.ProducerLinger),
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create franz-go client for producer: %w", err)
	}

	return cl, nil
}

func parseAcks(raw string) kgo.Acks {
	switch strings.TrimSpace(strings.ToLower(raw)) {
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
