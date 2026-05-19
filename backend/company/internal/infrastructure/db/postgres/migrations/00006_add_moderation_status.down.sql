ALTER TABLE vacancies DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE companies DROP COLUMN IF EXISTS moderation_status;

DROP TYPE IF EXISTS moderation_status;
