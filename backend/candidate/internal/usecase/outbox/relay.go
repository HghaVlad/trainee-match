package outbox

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
)

type RelayRepository interface {
	ListPendingAndSetProcessing(ctx context.Context, limit int, workerNumber, totalWorkers int) ([]Message, error)
	Save(ctx context.Context, msgs []Message) error
	ResetStaleProcessing(ctx context.Context, staleTimeout time.Duration) error
}

type Relay struct {
	repo      RelayRepository
	producer  Producer
	cfg       config.Outbox
	logger    *slog.Logger
	trManager *manager.Manager
}

func NewRelay(
	repo RelayRepository,
	producer Producer,
	config config.Outbox,
	logger *slog.Logger,
	trmanager *manager.Manager,
) *Relay {
	return &Relay{repo: repo, producer: producer, cfg: config, logger: logger, trManager: trmanager}
}

func (r *Relay) Run(ctx context.Context) {

	wg := sync.WaitGroup{}
	for i := range r.cfg.RelayWorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// WorkerNumber for having the same AggregateID only in one worker (select with hash % WorkerCount)
			r.runWorker(ctx, i)
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.runResetStaleProcessor(ctx)
	}()
	wg.Wait()
}

func (r *Relay) runWorker(ctx context.Context, workerNumber int) {

	sleep := r.cfg.RelayMinSleep

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processed := r.process(ctx, workerNumber)

		if processed != 0 {
			sleep = r.cfg.RelayMinSleep
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleep):
		}

		if sleep < r.cfg.RelayMaxSleep {
			sleep = min(sleep*2, r.cfg.RelayMaxSleep)
		}

	}
}

func (r *Relay) process(ctx context.Context, workerNumber int) int {

	var messages []Message
	var err error
	err = r.trManager.Do(ctx, func(ctx context.Context) error {
		messages, err = r.repo.ListPendingAndSetProcessing(
			ctx,
			r.cfg.RelayBatchSize,
			workerNumber,
			r.cfg.RelayWorkerCount,
		)
		return err
	})
	if err != nil {
		r.logger.WarnContext(ctx, "Error with getting messages from repository", slog.String("reason", err.Error()))
		return 0
	}

	results := r.producer.ProduceOutBox(ctx, messages)
	for i, result := range results {
		r.updateMessageWithResult(&messages[i], result)
	}

	err = r.trManager.Do(ctx, func(ctx context.Context) error {
		err = r.repo.Save(ctx, messages)
		return err
	})
	if err != nil {
		r.logger.WarnContext(ctx, "Error saving messages", slog.String("reason", err.Error()))
	}

	return len(results)

}

func (r *Relay) updateMessageWithResult(message *Message, result ProduceResult) {
	message.AttemptCount++
	if result.Err == nil {
		message.Status = StatusSent
		message.SentAt = result.SentAt
		return
	}

	err := result.Err.Error()
	message.LastError = &err
	if message.AttemptCount >= r.cfg.MaxAttempts || result.Unretryable {
		message.Status = StatusFailed
		return
	}

	message.Status = StatusPending
	delay := r.cfg.BaseRetryDelay * (1 << (message.AttemptCount)) // i. e. 5 10 20 40 80 160 ...
	message.NextAttemptAt = time.Now().UTC().Add(delay)
}

// To reset failed messages with left status = processing
func (r *Relay) runResetStaleProcessor(ctx context.Context) {
	if err := r.repo.ResetStaleProcessing(ctx, r.cfg.ResetStaleTime); err != nil {
		r.logger.WarnContext(ctx, "initial reset stale failed", "error", err)
	}

	t := time.NewTicker(r.cfg.ResetStaleTime)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := r.repo.ResetStaleProcessing(ctx, r.cfg.ResetStaleTime)
			if err != nil {
				r.logger.WarnContext(
					ctx,
					"Error resetting stale processor after timeout",
					slog.String("reason", err.Error()),
				)
			}
		}

	}
}
