package kafka

import (
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
)

func NewClient(config config.Kafka) (*kgo.Client, error) {

	opts := []kgo.Opt{
		kgo.SeedBrokers(config.Brokers...),
		kgo.ClientID(config.ClientID),
		kgo.RequiredAcks(parseAcks(config.ProducerAcks)),
		kgo.ProducerLinger(config.ProducerLinger),
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),
	}

	client, err := kgo.NewClient(opts...)
	return client, err
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
