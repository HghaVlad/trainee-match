DROP TRIGGER IF EXISTS update_resumes_updated_at ON resumes;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP INDEX IF EXISTS idx_resumes_candidate_id;
DROP TABLE IF EXISTS resumes;