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

CREATE TYPE vacancy_status_enum AS ENUM (
    'draft',
    'published',
    'archived'
);


CREATE TABLE vacancies (
    id                          UUID PRIMARY KEY,

    company_id                  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    created_by_user_id          UUID NOT NULL,

    title                       TEXT NOT NULL,
    description                 TEXT NOT NULL,

    work_format                 work_format_enum NOT NULL,
    city                        TEXT,

    duration_from_days          INT CHECK (duration_from_days > 0),
    duration_to_days            INT CHECK (duration_to_days > 0),

    employment_type             employment_type_enum NOT NULL,
    hours_per_week_from         INT CHECK (hours_per_week_from > 0),
    hours_per_week_to           INT CHECK (hours_per_week_to > 0),

    flexible_schedule           BOOLEAN NOT NULL DEFAULT false,

    is_paid                     BOOLEAN NOT NULL DEFAULT false,
    salary_from                 INT CHECK (salary_from >= 0),
    salary_to                   INT CHECK (salary_to >= 0),

    internship_to_offer         BOOLEAN NOT NULL DEFAULT false,

    status                      vacancy_status_enum NOT NULL DEFAULT 'draft',

    published_at                TIMESTAMPTZ,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- duration_from <= duration_to
ALTER TABLE vacancies
ADD CONSTRAINT chk_duration_range
CHECK (
    duration_from_days IS NULL
    OR duration_to_days IS NULL
    OR duration_from_days <= duration_to_days
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

-- if unpaid, no salaries
ALTER TABLE vacancies
ADD CONSTRAINT chk_salary_paid_logic
CHECK (
    (is_paid = false AND salary_from IS NULL AND salary_to IS NULL)
    OR
    (is_paid = true)
);

-- published -> published_at != null
ALTER TABLE vacancies
ADD CONSTRAINT chk_published_logic
CHECK (
    (status = 'published' AND published_at IS NOT NULL)
        OR
    (status != 'published' AND published_at IS NULL)
);

-- for published vacancy listing
CREATE INDEX idx_vacancies_feed
    ON vacancies (published_at DESC, id DESC)
    WHERE status = 'published';

CREATE INDEX idx_vacancies_company_feed
    ON vacancies (company_id, published_at DESC, id DESC)
    WHERE status = 'published';

CREATE INDEX idx_vacancies_salary_feed
    ON vacancies(salary_from DESC NULLS LAST, salary_to DESC NULLS LAST, id DESC)
    WHERE status = 'published';

CREATE INDEX idx_vacancies_salary_feed_asc
    ON vacancies(salary_from ASC NULLS LAST, salary_to ASC NULLS LAST, id ASC)
    WHERE status = 'published';

CREATE TRIGGER trg_vacancy_updated_at
BEFORE UPDATE ON vacancies
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
