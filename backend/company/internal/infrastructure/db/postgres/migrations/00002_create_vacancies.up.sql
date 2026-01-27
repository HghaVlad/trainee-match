CREATE TYPE work_format_enum AS ENUM (
    'onsite',
    'remote',
    'hybrid'
);

CREATE TYPE employment_type_enum AS ENUM (
    'full_time',
    'part_time',
    'internship'
);


CREATE TABLE vacancies (
    id                          UUID PRIMARY KEY,
    company_id                  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,

    title                       TEXT NOT NULL,
    description                 TEXT NOT NULL,

    work_format                 work_format_enum NOT NULL,
    city                        TEXT,

    duration_from_months        INT CHECK (duration_from_months > 0),
    duration_to_months          INT CHECK (duration_to_months > 0),

    employment_type             employment_type_enum NOT NULL,
    hours_per_week_from         INT CHECK (hours_per_week_from > 0),
    hours_per_week_to           INT CHECK (hours_per_week_to > 0),

    flexible_schedule           BOOLEAN NOT NULL DEFAULT false,

    is_paid                     BOOLEAN NOT NULL DEFAULT false,
    salary_from                 INT CHECK (salary_from >= 0),
    salary_to                   INT CHECK (salary_to >= 0),

    internship_to_offer         BOOLEAN NOT NULL DEFAULT false,

    is_active                   BOOLEAN NOT NULL DEFAULT true,
    published_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- duration_from <= duration_to
ALTER TABLE vacancies
ADD CONSTRAINT chk_duration_range
CHECK (
    duration_from_months IS NULL
    OR duration_to_months IS NULL
    OR duration_from_months <= duration_to_months
);

-- hours_from <= hours_to
ALTER TABLE vacancies
ADD CONSTRAINT chk_hours_range
CHECK (
    hours_per_week_from IS NULL
    OR hours_per_week_to IS NULL
    OR hours_per_week_from <= hours_per_week_to
);

-- salary_from <= salary_to
ALTER TABLE vacancies
ADD CONSTRAINT chk_salary_range
CHECK (
    salary_from IS NULL
    OR salary_to IS NULL
    OR salary_from <= salary_to
);


CREATE INDEX idx_vacancies_active
    ON vacancies (is_active)
    WHERE is_active = true;

CREATE INDEX idx_vacancies_work_format
    ON vacancies (work_format);

CREATE INDEX idx_vacancies_employment_type
    ON vacancies (employment_type);

CREATE INDEX idx_vacancies_company_id
    ON vacancies (company_id);

CREATE INDEX idx_vacancies_published_at
    ON vacancies (published_at DESC);


CREATE TRIGGER trg_vacancy_updated_at
BEFORE UPDATE ON vacancies
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
