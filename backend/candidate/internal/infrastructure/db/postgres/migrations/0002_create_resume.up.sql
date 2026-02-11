CREATE
TABLE IF NOT EXISTS resumes
(
    id           UUID DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    candidate_id UUID                           NOT NULL,
    name         VARCHAR(255)                   NOT NULL,
    status       INTEGER                        NOT NULL DEFAULT 0,
    data         JSONB                          NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_resumes_candidate_id ON resumes(candidate_id);
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

CREATE TRIGGER update_resumes_updated_at
    BEFORE UPDATE
    ON resumes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();