CREATE TABLE IF NOT EXISTS skills
(
    id         UUID DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name       VARCHAR(100)                   NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_skills_updated_at
    BEFORE UPDATE ON skills
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();