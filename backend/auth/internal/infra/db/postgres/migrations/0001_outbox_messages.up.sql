CREATE TABLE IF NOT EXISTS outbox_messages (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    aggregate_seq BIGINT NOT NULL DEFAULT 0,
    topic TEXT NOT NULL,
    key BYTEA NOT NULL,
    payload BYTEA NOT NULL,
    headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    schema_id INT NOT NULL,
    event_type TEXT NOT NULL,
    status TEXT NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ,
    last_error TEXT,
    next_attempt_at TIMESTAMPTZ NOT NULL,
    failed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_outbox_messages_status_next_attempt
    ON outbox_messages (status, next_attempt_at, created_at);

