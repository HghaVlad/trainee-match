CREATE TYPE outbox_status AS ENUM ('pending', 'processing', 'sent', 'failed');

CREATE TABLE outbox
(
    id              UUID PRIMARY KEY,
    aggregate_id    UUID          NOT NULL,
    aggregate_seq   BIGINT        NOT NULL,
    topic           TEXT          NOT NULL,
    event_type      TEXT          NOT NULL,
    schema_id       INT           NOT NULL,
    headers         JSONB         NOT NULL DEFAULT '{}'::jsonb,
    key             BYTEA         NOT NULL,
    payload         BYTEA         NOT NULL,
    status          outbox_status NOT NULL DEFAULT 'pending',
    attempt_count   BIGINT        NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    sent_at         TIMESTAMPTZ,
    last_error      TEXT,
    next_attempt_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    failed_at       TIMESTAMPTZ,
    started_at      TIMESTAMPTZ
);

CREATE TABLE outbox_seq
(
    aggregate_id UUID PRIMARY KEY,
    current_seq  BIGINT NOT NULL DEFAULT 0
);

-- Индекс для эффективного захвата pending‑сообщений relay‑воркерами
-- Позволяет быстро выполнять запрос:
--   SELECT * FROM outbox
--   WHERE status = 'pending'
--     AND next_attempt_at <= now()
--     AND abs(hashtext(aggregate_id::text)) % $2 = $3
CREATE INDEX idx_outbox_pending_claim
    ON outbox (abs(hashtext(aggregate_id::text)), next_attempt_at, id)
    WHERE status = 'pending';