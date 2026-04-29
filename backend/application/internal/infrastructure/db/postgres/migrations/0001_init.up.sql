CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE application_status_enum AS ENUM (
    'submitted',
    'seen',
    'interview',
    'rejected',
    'offer',
    'withdrawn'
    );

CREATE TYPE application_actor_enum AS ENUM (
    'candidate',
    'hr',
    'system'
    );

CREATE TYPE vacancy_status_enum AS ENUM (
    'published',
    'archived'
    );

CREATE type resume_status_enum AS ENUM (
    'draft',
    'published'
    );

CREATE TABLE application_snapshots
(
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),

    resume_id    UUID        NOT NULL,
    candidate_id UUID        NOT NULL,

    -- snapshot payload
    resume_name  TEXT,
    resume_data  JSONB       NOT NULL,

    full_name    TEXT        NOT NULL,
    email        TEXT        NOT NULL,
    telegram     TEXT,

    hash         TEXT        NOT NULL UNIQUE, -- deduplication

    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_snapshots_candidate ON application_snapshots (candidate_id);
CREATE INDEX idx_snapshots_resume ON application_snapshots (resume_id);


CREATE TABLE applications
(
    id           UUID PRIMARY KEY                 DEFAULT gen_random_uuid(),

    resume_id    UUID                    NOT NULL,
    candidate_id UUID                    NOT NULL,
    vacancy_id   UUID                    NOT NULL,
    company_id   UUID                    NOT NULL,

    snapshot_id  UUID                    NOT NULL REFERENCES application_snapshots (id),

    status       application_status_enum NOT NULL,
    cover_letter TEXT,

    created_at   TIMESTAMPTZ             NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ             NOT NULL DEFAULT now()
);

-- invariant: only one active application
CREATE UNIQUE INDEX uniq_active_application
    ON applications (candidate_id, vacancy_id)
    WHERE status IN ('submitted', 'seen', 'interview');

-- hr get by company
CREATE INDEX idx_applications_company ON applications (company_id);

-- candidate get mine
CREATE INDEX idx_applications_candidate ON applications (candidate_id);

-- hr get by vacancy
CREATE INDEX idx_applications_vacancy ON applications (vacancy_id);


CREATE TABLE application_status_history
(
    id                 UUID PRIMARY KEY                 DEFAULT gen_random_uuid(),

    application_id     UUID                    NOT NULL REFERENCES applications (id) ON DELETE CASCADE,

    status             application_status_enum NOT NULL,

    changed_by_user_id UUID,
    changed_by_role    application_actor_enum  NOT NULL,

    comment            TEXT,

    created_at         TIMESTAMPTZ             NOT NULL DEFAULT now()
);

CREATE INDEX idx_status_history_app ON application_status_history (application_id);


CREATE TABLE candidate_projection
(
    id         UUID PRIMARY KEY,

    full_name  TEXT,
    email      TEXT,
    telegram   TEXT,

    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);


CREATE TABLE resume_projection
(
    id           UUID PRIMARY KEY,

    candidate_id UUID               NOT NULL,

    name         TEXT               NOT NULL,
    data         JSONB              NOT NULL,
    status       resume_status_enum NOT NULL,

    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ
);

CREATE INDEX idx_resume_candidate ON resume_projection (candidate_id);


CREATE TABLE vacancy_projection
(
    id           UUID PRIMARY KEY,

    company_id   UUID                NOT NULL,
    company_name TEXT                NOT NULL,

    title        TEXT                NOT NULL,

    status       vacancy_status_enum NOT NULL,

    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ
);

CREATE INDEX idx_vacancies_company ON vacancy_projection (company_id);


CREATE TABLE company_members
(
    user_id    UUID NOT NULL,
    company_id UUID NOT NULL,
    role       TEXT,

    PRIMARY KEY (user_id, company_id)
);

CREATE INDEX idx_company_members_user ON company_members (user_id);

-- company_id in applications/company_members has no local компаний-таблицы;
-- FK добавлять не к чему, это остаётся на уровне интеграции с company-сервисом.
