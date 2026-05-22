CREATE TYPE moderation_status AS ENUM (
    'ok',
    'hidden'
);

ALTER TABLE vacancy_projection
    ADD COLUMN moderation_status moderation_status NOT NULL DEFAULT 'ok';

ALTER TABLE vacancy_projection
    ADD COLUMN company_moderation_status moderation_status NOT NULL DEFAULT 'ok';
