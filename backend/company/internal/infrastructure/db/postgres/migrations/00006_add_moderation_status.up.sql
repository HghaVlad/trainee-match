CREATE TYPE moderation_status AS ENUM (
    'ok',
    'hidden'
);

ALTER TABLE companies
    ADD COLUMN moderation_status moderation_status NOT NULL DEFAULT 'ok';

ALTER TABLE vacancies
    ADD COLUMN moderation_status moderation_status NOT NULL DEFAULT 'ok';
