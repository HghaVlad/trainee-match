package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/eventhandler"
)

type Consumer struct {
	cl      *kgo.Client
	handler *eventhandler.Handler

	workers map[tp]*partitionWorker
	mu      sync.RWMutex

	logger *slog.Logger
}

func NewConsumer(cfg config.Kafka, handler *eventhandler.Handler, logger *slog.Logger) (*Consumer, error) {
	consumer := &Consumer{
		handler: handler,
		workers: make(map[tp]*partitionWorker, partitionsCountExpectedUnder),
		mu:      sync.RWMutex{},
		logger:  logger,
	}

	cl, err := NewClientForConsumer(cfg, consumer)
	if err != nil {
		return nil, fmt.Errorf("new kafka client: %w", err)
	}

	consumer.cl = cl
	return consumer, nil
}

func (c *Consumer) Poll(ctx context.Context) {
	go func() {
		<-ctx.Done()
		c.cl.Close() // unblocks PollFetches, triggers shutdown
	}()

	for {
		fetches := c.cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			c.logger.WarnContext(ctx, "kafka consume fetch fail", "errors", errs)
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			w := c.getOrCreateWorker(ctx, p.Topic, p.Partition)

			if w.stopped.Load() {
				return
			}

			for _, r := range p.Records {
				w.jobs <- r
			}
		})

		c.logger.InfoContext(ctx, "got kafka fetches", "cnt", len(fetches))

		select {
		case <-ctx.Done():
			c.shutdown(ctx)
			return
		default:
		}
	}
}

func (c *Consumer) handle(ctx context.Context, record *kgo.Record) {
	msg := &eventhandler.Event{
		Payload: record.Value,
		Key:     record.Key,
		Topic:   record.Topic,
		Headers: headersMap(record),
	}

	c.handler.HandleMsg(ctx, msg)
}

// for each of assigned partitions create worker and run it
func (c *Consumer) onAssigned(ctx context.Context, partitions map[string][]int32) {
	for topic := range partitions {
		for _, p := range partitions[topic] {
			c.getOrCreateWorker(ctx, topic, p)
		}
	}
}

// delete workers of revoked partitions from map,
// each of them stops receiving work and handles the rest of received work, commits offset
func (c *Consumer) onRevoked(_ context.Context, partitions map[string][]int32) {
	var toStop []*partitionWorker

	c.mu.Lock()

	for topic := range partitions {
		for _, p := range partitions[topic] {
			key := tp{topic: topic, partition: p}

			w, exists := c.workers[key]
			if exists {
				toStop = append(toStop, w)
				delete(c.workers, key)
			}
		}
	}

	c.mu.Unlock()

	wg := &sync.WaitGroup{}

	for _, w := range toStop {
		wg.Go(w.shutdown)
	}

	wg.Wait()
}

// handles the rest of received records
func (c *Consumer) shutdown(ctx context.Context) {
	c.logger.InfoContext(ctx, "kafka consumer gracefully shutting down")

	wg := &sync.WaitGroup{}

	for _, w := range c.workers {
		wg.Go(w.shutdown)
	}

	wg.Wait()
}

func (c *Consumer) commitAsync(ctx context.Context, topic string, partition int32, offset int64) {
	c.cl.CommitOffsets(ctx,
		map[string]map[int32]kgo.EpochOffset{
			topic: {partition: {Offset: offset + 1}},
		},
		func(_ *kgo.Client, _ *kmsg.OffsetCommitRequest, _ *kmsg.OffsetCommitResponse, err error) {
			if err != nil {
				c.logger.Warn("kafka async commit offset fail",
					"topic", topic, "partition", partition,
					"offset", offset, "error", err)
			}
		})
}

func (c *Consumer) commitSync(ctx context.Context, topic string, partition int32, offset int64) {
	c.cl.CommitOffsetsSync(ctx,
		map[string]map[int32]kgo.EpochOffset{
			topic: {partition: {Offset: offset + 1}},
		},
		func(_ *kgo.Client, _ *kmsg.OffsetCommitRequest, _ *kmsg.OffsetCommitResponse, err error) {
			if err != nil {
				c.logger.Warn("kafka commit offset fail",
					"topic", topic, "partition", partition,
					"offset", offset, "error", err)
			}
		})
}

// gets or creates and runs a new worker
func (c *Consumer) getOrCreateWorker(ctx context.Context, topic string, partition int32) *partitionWorker {
	c.mu.RLock()
	key := tp{topic: topic, partition: partition}
	if w, ok := c.workers[key]; ok {
		c.mu.RUnlock()
		return w
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// double check
	if w, ok := c.workers[key]; ok {
		return w
	}

	w := &partitionWorker{
		jobs:       make(chan *kgo.Record, partitionWorkerJobsCap),
		consumer:   c,
		topic:      topic,
		partition:  partition,
		lastOffset: -1, // -1 won't be commited
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
		logger:     c.logger,
	}
	c.workers[key] = w

	go w.run(ctx)

	return w
}

func headersMap(record *kgo.Record) map[string][]byte {
	m := make(map[string][]byte, len(record.Headers))

	for _, kv := range record.Headers {
		m[kv.Key] = kv.Value
	}

	return m
}
