package kafka

import (
	"context"
	"log/slog"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/eventhandler"
)

type Consumer struct {
	mu        sync.Mutex
	consumers map[string]map[int32]*PartitionConsumer
	Client    *kgo.Client
	handler   *eventhandler.Handler
}

func NewConsumer(handler *eventhandler.Handler) *Consumer {
	return &Consumer{consumers: make(map[string]map[int32]*PartitionConsumer), handler: handler}
}

func (c *Consumer) Assigned(_ context.Context, _ *kgo.Client, assigned map[string][]int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for topic, partitions := range assigned {
		if _, ok := c.consumers[topic]; !ok {
			c.consumers[topic] = make(map[int32]*PartitionConsumer)
		}
		for _, partition := range partitions {
			if _, ok := c.consumers[topic][partition]; ok {
				continue
			}
			pr := NewPartitionConsumer(topic, partition, c.Client, c.handler)
			c.consumers[topic][partition] = pr

			go pr.ConsumePartition(topic, partition)
		}
	}
}

func (c *Consumer) Revoked(_ context.Context, _ *kgo.Client, revoked map[string][]int32) {
	c.mu.Lock()
	toStop := make([]*PartitionConsumer, 0)
	for topic, partitions := range revoked {
		if _, ok := c.consumers[topic]; !ok {
			continue
		}
		for _, partition := range partitions {
			pr := c.consumers[topic][partition]
			delete(c.consumers[topic], partition)
			if len(c.consumers[topic]) == 0 {
				delete(c.consumers, topic)
			}
			toStop = append(toStop, pr)
		}
	}
	c.mu.Unlock()

	wg := sync.WaitGroup{}
	for _, pr := range toStop {
		wg.Go(pr.Shutdown)
	}
	wg.Wait()
}

func (c *Consumer) Poll(ctx context.Context) {
	if c.Client == nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			c.Shutdown()
			return
		default:
		}

		fetches := c.Client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			c.Shutdown()
			return
		}

		if err := fetches.Err(); err != nil {
			slog.Warn("kafka consumer has fetches error", "error", err)
		}

		fetches.EachTopic(func(partitionTopic kgo.FetchTopic) {
			c.mu.Lock()
			topicConsumers := c.consumers[partitionTopic.Topic]
			c.mu.Unlock()
			if topicConsumers == nil {
				return
			}
			partitionTopic.EachPartition(func(partition kgo.FetchPartition) {
				pr, ok := topicConsumers[partition.Partition]
				if !ok {
					return
				}
				select {
				case <-pr.quit:
				case pr.records <- partition.Records:
				}
			})
		})
	}
}

func (c *Consumer) Shutdown() {
	c.mu.Lock()
	toStop := make([]*PartitionConsumer, 0)
	for topic, partitions := range c.consumers {
		if _, ok := c.consumers[topic]; !ok {
			continue
		}
		for i, pr := range partitions {
			delete(c.consumers[topic], i)
			if len(c.consumers[topic]) == 0 {
				delete(c.consumers, topic)
			}
			toStop = append(toStop, pr)
		}
	}
	c.mu.Unlock()

	wg := sync.WaitGroup{}
	for _, pr := range toStop {
		wg.Go(pr.Shutdown)
	}
	wg.Wait()
}
