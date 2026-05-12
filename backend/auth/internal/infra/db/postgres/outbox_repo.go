package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/services/outbox"
)

const (
	statusPending    = "pending"
	statusProcessing = "processing"
	statusSent       = "sent"
	statusFailed     = "failed"
)

type OutboxRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewOutboxRepo(pool *pgxpool.Pool, getter *trmpgx.CtxGetter) *OutboxRepo {
	return &OutboxRepo{pool: pool, getter: getter}
}

func (r *OutboxRepo) Create(ctx context.Context, msg outbox.Message) error {
	headers, err := json.Marshal(msg.Headers)
	if err != nil {
		return fmt.Errorf("marshal headers: %w", err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)

	_, err = conn.Exec(ctx, `
		INSERT INTO outbox_messages (
			id,
			aggregate_id,
			aggregate_seq,
			topic,
			key,
			payload,
			headers,
			schema_id,
			event_type,
			status,
			attempt_count,
			max_attempts,
			created_at,
			next_attempt_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14
		)
	`,
		msg.ID,
		msg.AggregateID,
		msg.AggregateSeq,
		msg.Topic,
		msg.Key,
		msg.Payload,
		headers,
		msg.SchemaID,
		string(msg.EventType),
		string(msg.Status),
		msg.AttemptCount,
		msg.MaxAttempts,
		msg.CreatedAt,
		msg.NextAttemptAt,
	)
	if err != nil {
		return fmt.Errorf("insert outbox message: %w", err)
	}

	return nil
}

func (r *OutboxRepo) ListPending(ctx context.Context, batchSize int) ([]outbox.Message, error) {

	conn := r.getter.DefaultTrOrDB(ctx, r.pool)

	rows, err := conn.Query(ctx, `
		WITH picked AS (
			SELECT id
			FROM outbox_messages
			WHERE status = $1
				AND next_attempt_at <= now()
			ORDER BY created_at
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)
		UPDATE outbox_messages
		SET status = $3
		WHERE id IN (SELECT id FROM picked)
		RETURNING id, aggregate_id, aggregate_seq, topic, key, payload, headers,
			schema_id, event_type, status, attempt_count, max_attempts,
			created_at, sent_at, last_error, next_attempt_at, failed_at
	`, statusPending, batchSize, statusProcessing)
	if err != nil {
		return nil, fmt.Errorf("list pending outbox: %w", err)
	}
	defer rows.Close()

	result := make([]outbox.Message, 0, batchSize)

	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list pending outbox rows: %w", err)
	}

	return result, nil
}

func (r *OutboxRepo) Save(ctx context.Context, msgs []outbox.Message) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.pool)

	batch := &pgx.Batch{}

	for i := range msgs {
		batch.Queue(`
			UPDATE outbox_messages
			SET status = $2,
				attempt_count = $3,
				sent_at = $4,
				last_error = $5,
				next_attempt_at = $6,
				failed_at = $7
			WHERE id = $1
		`,
			msgs[i].ID,
			string(msgs[i].Status),
			msgs[i].AttemptCount,
			msgs[i].SentAt,
			msgs[i].LastError,
			msgs[i].NextAttemptAt,
			msgs[i].FailedAt,
		)
	}

	br := conn.SendBatch(ctx, batch)
	defer br.Close()

	for range msgs {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("save outbox messages: %w", err)
		}
	}

	return nil
}

func scanMessage(row pgx.Row) (outbox.Message, error) {
	var msg outbox.Message
	var headersBytes []byte
	var eventType string
	var status string
	var sentAt *time.Time
	var lastError *string
	var failedAt *time.Time

	if err := row.Scan(
		&msg.ID,
		&msg.AggregateID,
		&msg.AggregateSeq,
		&msg.Topic,
		&msg.Key,
		&msg.Payload,
		&headersBytes,
		&msg.SchemaID,
		&eventType,
		&status,
		&msg.AttemptCount,
		&msg.MaxAttempts,
		&msg.CreatedAt,
		&sentAt,
		&lastError,
		&msg.NextAttemptAt,
		&failedAt,
	); err != nil {
		return outbox.Message{}, fmt.Errorf("scan outbox message: %w", err)
	}

	headers := make(map[string]string)
	if len(headersBytes) > 0 {
		if err := json.Unmarshal(headersBytes, &headers); err != nil {
			return outbox.Message{}, fmt.Errorf("unmarshal headers: %w", err)
		}
	}

	msg.Headers = headers
	msg.EventType = outbox.EventType(eventType)
	msg.Status = outbox.Status(status)
	msg.SentAt = sentAt
	msg.LastError = lastError
	msg.FailedAt = failedAt

	return msg, nil
}
