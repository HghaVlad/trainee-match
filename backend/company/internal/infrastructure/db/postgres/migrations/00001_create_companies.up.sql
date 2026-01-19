CREATE TABLE companies (
   id          UUID PRIMARY KEY,
   name        VARCHAR(128) UNIQUE NOT NULL,
   description TEXT,
   website     TEXT,
   logo_key    TEXT,
   owner_id    UUID NOT NULL,

   created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);


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


