package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
)

type OutboxRepo struct {
	db *pgxpool.Pool
}

func NewOutboxRepo(db *pgxpool.Pool) *OutboxRepo {
	return &OutboxRepo{
		db: db,
	}
}

func (r *OutboxRepo) Create(ctx context.Context, msg outbox.Message) error {
	q := postgres.GetQuerier(ctx, r.db)

	headersB, err := json.Marshal(msg.Headers)
	if err != nil {
		return fmt.Errorf("create outbox msg: marshal headers: %w", err)
	}

	const query = `INSERT INTO outbox
    (id, topic, key, payload, headers, schema_id,
     event_type, status, max_attempts, next_attempt_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = q.Exec(ctx, query,
		msg.ID, msg.Topic, msg.Key, msg.Payload, headersB, msg.SchemaID,
		msg.EventType, msg.Status, msg.MaxAttempts, msg.NextAttemptAt, msg.CreatedAt)

	if err != nil {
		return fmt.Errorf("create outbox msg: %w", err)
	}

	return nil
}

// ListPending uses skip locked for better throughput,
// Messages are ordered by (next_attempt_at, attempt_count),
// so that messages with lesser attempts had higher priority
func (r *OutboxRepo) ListPending(ctx context.Context, limit int) ([]outbox.Message, error) {
	q := postgres.GetQuerier(ctx, r.db)

	const query = `SELECT
		id, topic, key, payload, headers, schema_id,
		event_type, status, attempt_count, max_attempts,
		next_attempt_at, created_at
		FROM outbox
		WHERE status = 'pending' AND next_attempt_at <= now()
		ORDER BY next_attempt_at, attempt_count > 0, attempt_count, id
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := q.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list outbox pending: %w", err)
	}

	return scanMessages(rows, scanPending)
}

func (r *OutboxRepo) Save(ctx context.Context, msgs []outbox.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	q := postgres.GetQuerier(ctx, r.db)

	const query = `UPDATE outbox 
		SET status = $1, sent_at = $2, attempt_count = $3, last_error = $4, next_attempt_at = $5, failed_at = $6
		WHERE id = $7`

	batch := &pgx.Batch{}

	for _, msg := range msgs {
		batch.Queue(
			query,
			msg.Status,
			msg.SentAt,
			msg.AttemptCount,
			msg.LastError,
			msg.NextAttemptAt,
			msg.FailedAt,
			msg.ID,
		)
	}

	br := q.SendBatch(ctx, batch)

	if err := br.Close(); err != nil {
		return fmt.Errorf("save outbox: %w", err)
	}

	return nil
}

func scanMessages(rows pgx.Rows,
	scanFunc func(rows pgx.Rows, msg *outbox.Message, headersB *[]byte) error,
) ([]outbox.Message, error) {
	defer rows.Close()
	var msgs []outbox.Message

	for rows.Next() {
		msg := outbox.Message{
			Headers: make(map[string]string),
		}

		headersB := make([]byte, 0)

		err := scanFunc(rows, &msg, &headersB)

		if err != nil {
			return nil, fmt.Errorf("scan outbox msgs: %w", err)
		}

		if err := json.Unmarshal(headersB, &msg.Headers); err != nil {
			return nil, fmt.Errorf("list outbox msgs: unmarshal headers: %w", err)
		}

		msgs = append(msgs, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan outbox msgs: %w", err)
	}

	return msgs, nil
}

func scanPending(rows pgx.Rows, msg *outbox.Message, headersB *[]byte) error {
	return rows.Scan(
		&msg.ID, &msg.Topic, &msg.Key, &msg.Payload, headersB, &msg.SchemaID,
		&msg.EventType, &msg.Status, &msg.AttemptCount, &msg.MaxAttempts,
		&msg.NextAttemptAt, &msg.CreatedAt,
	)
}
