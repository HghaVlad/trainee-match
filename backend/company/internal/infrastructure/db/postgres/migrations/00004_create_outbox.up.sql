CREATE TYPE outbox_message_status AS ENUM (
    'pending',
    'sent',
    'failed'
);

CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    seq INT NOT NULL,
    status outbox_message_status NOT NULL DEFAULT 'pending',
    key BYTEA NOT NULL,
    payload BYTEA NOT NULL,
    headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    schema_id INT NOT NULL,
    topic TEXT NOT NULL,
    event_type TEXT NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    max_attempts INT NOT NULL DEFAULT 3 CHECK (max_attempts >= 0),
    last_error TEXT,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,

    CONSTRAINT sent_requires_sent_at CHECK (
            (status = 'sent' AND sent_at IS NOT NULL) OR
            (status <> 'sent' AND sent_at IS NULL)),

    CONSTRAINT dead_requires_failed_at CHECK (
        (status = 'failed' AND failed_at IS NOT NULL) OR
        (status <> 'failed' AND failed_at IS NULL))
);

CREATE INDEX idx_outbox_pending_poll
    ON outbox(next_attempt_at ASC, attempt_count ASC, id ASC)
    WHERE status = 'pending';

CREATE INDEX idx_outbox_pending_aggregate_id_seq
    ON outbox(aggregate_id, seq)
    WHERE status = 'pending';

CREATE TABLE outbox_seq (
    aggregate_id uuid PRIMARY KEY,
    seq INT NOT NULL
);
