package kafka

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/eventhandler"
)

type PartitionConsumer struct {
	topic     string
	partition int32
	quit      chan struct{}
	records   chan []*kgo.Record
	stopped   atomic.Bool
	done      chan struct{}
	client    *kgo.Client
	Handler   *eventhandler.Handler
	offset    int64
}

func NewPartitionConsumer(
	topic string,
	partition int32,
	client *kgo.Client,
	handler *eventhandler.Handler,
) *PartitionConsumer {
	return &PartitionConsumer{
		topic:     topic,
		partition: partition,
		quit:      make(chan struct{}),
		records:   make(chan []*kgo.Record, 100),
		done:      make(chan struct{}),
		client:    client,
		Handler:   handler,
		offset:    -1,
	}
}

func (c *PartitionConsumer) ConsumePartition(_ string, _ int32) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.quit:
			for {
				select {
				case records := <-c.records:
					for _, record := range records {
						c.handleRecord(record)
						c.offset = record.Offset
					}
					c.CommitSync()
				default:
					// All records handled returning
					c.CommitSync()
					close(c.done)
					return
				}
			}
		case records := <-c.records:
			for _, record := range records {
				c.handleRecord(record)
				c.offset = record.Offset
			}
		case <-ticker.C:
			c.Commit()
		}
	}
}

func (c *PartitionConsumer) handleRecord(record *kgo.Record) {
	newEvent := eventhandler.Event{
		Topic:   record.Topic,
		Key:     record.Key,
		Payload: record.Value,
		Headers: recordHeadersToMapHeaders(record.Headers),
	}
	c.Handler.HandleEvent(context.Background(), newEvent)
}

// Shutdown with handling all left records
func (c *PartitionConsumer) Shutdown() {
	if !c.stopped.Load() {
		c.stopped.Store(true)
		close(c.quit)
		<-c.done
	}
}

func (c *PartitionConsumer) Commit() {
	if c.offset == -1 {
		return
	}
	offsets := map[string]map[int32]kgo.EpochOffset{
		c.topic: {c.partition: {Offset: c.offset + 1}},
	}
	c.client.CommitOffsets(context.Background(), offsets,
		func(client *kgo.Client, request *kmsg.OffsetCommitRequest, response *kmsg.OffsetCommitResponse, err error) {
			if err != nil {
				slog.Warn("kafka commit offset failed", "topic", c.topic, "partition", c.partition, "err", err)
			}
		})
}

func (c *PartitionConsumer) CommitSync() {
	if c.offset == -1 {
		return
	}
	offsets := map[string]map[int32]kgo.EpochOffset{
		c.topic: {c.partition: {Offset: c.offset + 1}},
	}
	c.client.CommitOffsetsSync(context.Background(), offsets,
		func(client *kgo.Client, request *kmsg.OffsetCommitRequest, response *kmsg.OffsetCommitResponse, err error) {
			if err != nil {
				slog.Warn("kafka commit offset sync failed", "topic", c.topic, "partition", c.partition, "err", err)
			}
		})
}

func recordHeadersToMapHeaders(hds []kgo.RecordHeader) map[string][]byte {
	headers := make(map[string][]byte)
	for _, header := range hds {
		headers[header.Key] = header.Value
	}
	return headers
}
