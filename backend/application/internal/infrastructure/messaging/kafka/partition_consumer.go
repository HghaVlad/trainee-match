package kafka

import (
	"fmt"
	"sync/atomic"

	"github.com/twmb/franz-go/pkg/kgo"
)

type PartitionConsumer struct {
	quit    chan struct{}
	records chan []*kgo.Record
	stopped atomic.Bool
	done    chan struct{}
}

func NewPartitionConsumer() *PartitionConsumer {
	return &PartitionConsumer{
		quit:    make(chan struct{}),
		records: make(chan []*kgo.Record, 100),
		done:    make(chan struct{}),
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
	fmt.Println("got record", record)
}

// Shutdown with handling all left records
func (c *PartitionConsumer) Shutdown() {
	if !c.stopped.Load() {
		c.stopped.Store(true)
		close(c.quit)
		<-c.done
	}
}
