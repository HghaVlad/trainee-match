CREATE TABLE IF NOT EXISTS candidates
(
    id
    UUID
    DEFAULT
    gen_random_uuid
(
) NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    phone VARCHAR
(
    20
),
    telegram VARCHAR
(
    50
),
    city VARCHAR
(
    100
),
    birthday DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP
                         WITH TIME ZONE DEFAULT NOW()
    );


CREATE INDEX IF NOT EXISTS idx_candidates_user_id ON candidates(user_id);


CREATE
OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at
= NOW();
RETURN NEW;
END;
$$
language 'plpgsql';

DROP TRIGGER IF EXISTS update_candidates_updated_at ON candidates;
CREATE TRIGGER update_candidates_updated_at
    BEFORE UPDATE
    ON candidates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();