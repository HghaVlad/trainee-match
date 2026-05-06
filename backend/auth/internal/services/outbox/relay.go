package outbox

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
)

// Relay is a group of workers that take messages from outbox repo
// and call broker producer, manage the results.
type Relay struct {
	producer Producer
	repo     RelayRepo
	cfg      config.Outbox
	logger   *slog.Logger
}

type RelayRepo interface {
	ListPending(ctx context.Context, batchSize int) ([]Message, error)
	Save(ctx context.Context, msgs []Message) error
}

func NewRelay(producer Producer, repo RelayRepo, cfg config.Outbox, logger *slog.Logger) *Relay {
	return &Relay{producer: producer, repo: repo, cfg: cfg, logger: logger}
}

// Run launches workers, is a blocking operation, stopped by ctx.Done.
func (re *Relay) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	for range re.cfg.RelayWorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			re.runWorker(ctx)
		}()
	}

	wg.Wait()
}

func (re *Relay) Process(ctx context.Context) int {
	processed := 0

	msgs, err := re.repo.ListPending(ctx, re.cfg.BatchSize)
	if err != nil {
		re.logger.WarnContext(ctx, "outbox list pending fail", "err", err)
		return processed
	}

	if len(msgs) == 0 {
		return processed
	}

	results := re.producer.ProduceOutbox(ctx, msgs)
	re.updateMsgsFromResults(msgs, results)

	if err := re.repo.Save(ctx, msgs); err != nil {
		re.logger.WarnContext(ctx, "relay outbox messages save fail", "cnt", len(msgs), "err", err)
	}

	processed = len(msgs)
	return processed
}

func (re *Relay) runWorker(ctx context.Context) {
	sleep := re.cfg.RelayMinSleep

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processed := re.Process(ctx)

		// don't sleep if there's work
		if processed > 0 {
			sleep = re.cfg.RelayMinSleep
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleep):
		}

		// increase sleep
		if sleep < re.cfg.RelayMaxSleep {
			sleep = min(sleep*2, re.cfg.RelayMaxSleep)
		}
	}
}

func (re *Relay) updateMsgsFromResults(msgs []Message, results []ProduceResult) {
	for i := range msgs {
		msgs[i].AttemptCount++

		switch {
		case results[i].SentAt != nil:
			markSent(&msgs[i], results[i])
		case results[i].Unretryable || msgs[i].AttemptCount >= re.cfg.MaxRetries:
			markFailed(&msgs[i], results[i])
		default:
			re.markForRetry(&msgs[i])
			msgs[i].Status = StatusPending
		}
	}
}

func markSent(msg *Message, res ProduceResult) {
	msg.Status = StatusSent
	msg.SentAt = res.SentAt
}

func (re *Relay) markForRetry(msg *Message) {
	delay := re.cfg.BaseRetryDelay * (1 << (msg.AttemptCount - 1))
	msg.NextAttemptAt = time.Now().UTC().Add(delay)
}

func markFailed(msg *Message, res ProduceResult) {
	msg.Status = StatusFailed

	errMsg := res.Err.Error()
	msg.LastError = &errMsg

	now := time.Now().UTC()
	msg.FailedAt = &now
}
