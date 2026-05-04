package kafka

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	commitTickerDuration         = 100 * time.Millisecond
	partitionWorkerJobsCap       = 1000
	partitionsCountExpectedUnder = 200
)

type tp struct {
	topic     string
	partition int32
}

type partitionWorker struct {
	jobs       chan *kgo.Record
	consumer   *Consumer
	topic      string
	partition  int32
	lastOffset int64
	stop       chan struct{}
	done       chan struct{}
	stopped    atomic.Bool
}

func (pw *partitionWorker) run(ctx context.Context) {
	ticker := time.NewTicker(commitTickerDuration)
	defer ticker.Stop()

	for {
		select {
		case r := <-pw.jobs:
			pw.consumer.handle(ctx, r) // handle or write to dlq after retries
			pw.lastOffset = r.Offset

		case <-ticker.C: // commit every 100 ms
			pw.consumer.commitAsync(ctx, pw.topic, pw.partition, pw.lastOffset)

		case <-pw.stop:
			// graceful drain, once, after shutdown is called
			for {
				select {
				case r := <-pw.jobs:
					pw.consumer.handle(ctx, r)
					pw.lastOffset = r.Offset
				default:
					pw.consumer.commitSync(ctx, pw.topic, pw.partition, pw.lastOffset)
					close(pw.done)
					return
				}
			}
		}
	}
}

// blocks until no new records to handle
func (pw *partitionWorker) shutdown() {
	pw.stopped.Store(true)
	close(pw.stop)
	<-pw.done
}

