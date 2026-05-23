ALTER TABLE resumes
    DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE candidates
    DROP COLUMN IF EXISTS full_name;

