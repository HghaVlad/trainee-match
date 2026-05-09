DROP INDEX IF EXISTS idx_outbox_pending_claim;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS aggregate_sequences;
DROP TYPE IF EXISTS outbox_status;