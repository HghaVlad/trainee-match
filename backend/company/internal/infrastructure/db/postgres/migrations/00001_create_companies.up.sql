CREATE TABLE companies (
   id          UUID PRIMARY KEY,
   name        VARCHAR(128) UNIQUE NOT NULL,
   description TEXT,
   website     TEXT,
   logo_key    TEXT,
   open_vacancies_count INT NOT NULL DEFAULT 0,
   CHECK (open_vacancies_count >= 0),

   created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_companies_open_vacancies_desc
    ON companies (open_vacancies_count DESC, name ASC);

CREATE INDEX idx_companies_created_at_desc
    ON companies (created_at DESC, name ASC);

-- Func for auto update_at update
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Trigger for auto updated_at update
CREATE TRIGGER trg_company_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


