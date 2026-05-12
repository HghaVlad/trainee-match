package repository

import (
	"context"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/common/outbox"
)

type Outbox struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewOutbox(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *Outbox {
	return &Outbox{db: db, getter: getter}
}

func (r *Outbox) Create(ctx context.Context, msg outbox.Message) error {

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	const querySeqRow = `INSERT INTO outbox_seq (aggregate_id, current_seq) VALUES ($1, 1) ON CONFLICT (aggregate_id)
     DO UPDATE SET current_seq = outbox_seq.current_seq + 1 RETURNING current_seq`

	var seq int
	err := conn.QueryRow(ctx, querySeqRow, msg.AggregateID).Scan(&seq)
	if err != nil {
		return fmt.Errorf("create outbox: update seq: %w", err)
	}

	const query = `INSERT INTO outbox (
            id, aggregate_id, aggregate_seq, topic, event_type, schema_id,
            headers, key, payload, status, attempt_count, created_at, next_attempt_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, $11, $12, $13
        )`

	_, err = conn.Exec(ctx, query,
		msg.ID,
		msg.AggregateID,
		seq,
		msg.Topic,
		string(msg.EventType),
		msg.SchemaID,
		msg.Headers,
		msg.Key,
		msg.Payload,
		msg.Status,
		msg.AttemptCount,
		msg.CreatedAt,
		msg.NextAttemptAt,
	)
	if err != nil {
		return fmt.Errorf("create outbox: %w", err)
	}
	return nil
}

func (r *Outbox) ListPendingAndSetProcessing(
	ctx context.Context,
	limit int,
	workerNumber, totalWorkers int,
) ([]outbox.Message, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	rows, err := conn.Query(ctx,
		`SELECT id, aggregate_id, aggregate_seq, topic, event_type, schema_id,
		        headers, key, payload, status, attempt_count, created_at,
		        sent_at, last_error, next_attempt_at, failed_at
		 FROM outbox
		 WHERE status = 'pending'
		   AND next_attempt_at <= now()
		   AND abs(hashtext(aggregate_id::text)) % $2 = $3
		 ORDER BY aggregate_seq
		 LIMIT $1
		 FOR UPDATE SKIP LOCKED`,
		limit, totalWorkers, workerNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("query pending: %w", err)
	}
	defer rows.Close()
	var msgs []outbox.Message
	for rows.Next() {
		var m outbox.Message
		err = rows.Scan(
			&m.ID, &m.AggregateID, &m.AggregateSeq,
			&m.Topic, &m.EventType, &m.SchemaID,
			&m.Headers, &m.Key, &m.Payload,
			&m.Status, &m.AttemptCount, &m.CreatedAt,
			&m.SentAt, &m.LastError, &m.NextAttemptAt,
			&m.FailedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	if len(msgs) == 0 {
		return nil, nil
	}

	// Обновляем статус и время старта
	now := time.Now().UTC()
	ids := make([]uuid.UUID, len(msgs))
	for i, m := range msgs {
		ids[i] = m.ID
	}
	_, err = conn.Exec(ctx,
		`UPDATE outbox o
		 SET status = 'processing', started_at = $1
		 FROM UNNEST($2::uuid[]) AS t(id)
		 WHERE o.id = t.id`,
		now, ids,
	)
	if err != nil {
		return nil, fmt.Errorf("mark processing: %w", err)
	}

	// Возвращаем сообщения с обновлёнными полями
	for i := range msgs {
		msgs[i].Status = outbox.StatusProcessing
	}

	return msgs, nil

}
func (r *Outbox) Save(ctx context.Context, msgs []outbox.Message) error {

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	n := len(msgs)
	if n == 0 {
		return nil
	}

	ids := make([]uuid.UUID, n)
	statuses := make([]string, n)
	attemptCounts := make([]int64, n)
	sentAts := make([]*time.Time, n)
	lastErrors := make([]*string, n)
	nextAttemptAts := make([]time.Time, n)
	failedAts := make([]*time.Time, n)

	for i, m := range msgs {
		ids[i] = m.ID
		statuses[i] = string(m.Status)
		attemptCounts[i] = m.AttemptCount
		sentAts[i] = m.SentAt
		lastErrors[i] = m.LastError
		nextAttemptAts[i] = m.NextAttemptAt
		failedAts[i] = m.FailedAt
	}

	_, err := conn.Exec(ctx,
		`UPDATE outbox o SET
            status          = t.status,
            attempt_count   = t.attempt_count,
            sent_at         = t.sent_at,
            last_error      = t.last_error,
            next_attempt_at = t.next_attempt_at,
            failed_at       = t.failed_at,
            started_at      = NULL
         FROM UNNEST(
            $1::uuid[],
            $2::outbox_status[],
            $3::bigint[],
            $4::timestamptz[],
            $5::text[],
            $6::timestamptz[],
            $7::timestamptz[]
         ) AS t(id, status, attempt_count, sent_at, last_error, next_attempt_at, failed_at)
         WHERE o.id = t.id`,
		ids,
		statuses,
		attemptCounts,
		sentAts,
		lastErrors,
		nextAttemptAts,
		failedAts,
	)
	if err != nil {
		return fmt.Errorf("batch save: %w", err)
	}
	return nil
}

func (r *Outbox) ResetStaleProcessing(ctx context.Context, staleTimeout time.Duration) error {
	cutoff := time.Now().UTC().Add(-staleTimeout)

	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	_, err := conn.Exec(ctx,
		`UPDATE outbox
         SET status = 'pending',
             started_at = NULL,
             attempt_count = attempt_count + 1,
             next_attempt_at = now()
         WHERE status = 'processing'
           AND started_at < $1`,
		cutoff,
	)
	if err != nil {
		return fmt.Errorf("reset stale processing: %w", err)
	}
	return nil
}
