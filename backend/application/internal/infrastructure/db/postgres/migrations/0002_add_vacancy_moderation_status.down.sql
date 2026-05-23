ALTER TABLE vacancy_projection DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE vacancy_projection DROP COLUMN IF EXISTS company_moderation_status;

DROP TYPE IF EXISTS moderation_status;