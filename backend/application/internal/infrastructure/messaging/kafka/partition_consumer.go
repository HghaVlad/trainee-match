package kafka

import (
	"context"
	"sync/atomic"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/eventhandler"
)

type PartitionConsumer struct {
	quit    chan struct{}
	records chan []*kgo.Record
	stopped atomic.Bool
	done    chan struct{}
	Handler *eventhandler.Handler
}

func NewPartitionConsumer(handler *eventhandler.Handler) *PartitionConsumer {
	return &PartitionConsumer{
		quit:    make(chan struct{}),
		records: make(chan []*kgo.Record, 100),
		done:    make(chan struct{}),
		Handler: handler,
	}
}

func (c *PartitionConsumer) ConsumePartition(_ string, _ int32) {
	for {
		select {
		case <-c.quit:
			for {
				select {
				case records := <-c.records:
					for _, record := range records {
						c.handleRecord(record)
					}
				default:
					// All records handled returning
					close(c.done)
					return
				}
			}
		case records := <-c.records:
			for _, record := range records {
				c.handleRecord(record)
			}
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

func recordHeadersToMapHeaders(hds []kgo.RecordHeader) map[string][]byte {
	headers := make(map[string][]byte)
	for _, header := range hds {
		headers[header.Key] = header.Value
	}
	return headers
}
