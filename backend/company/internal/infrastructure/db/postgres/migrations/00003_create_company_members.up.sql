CREATE TYPE company_member_role_enum AS ENUM (
    'recruiter',
    'admin'
);

CREATE TABLE IF NOT EXISTS company_members (
    user_id UUID NOT NULL,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    role company_member_role_enum NOT NULL DEFAULT 'recruiter',

    PRIMARY KEY (user_id, company_id)
);

CREATE INDEX idx_company_members_company
    ON company_members(company_id);
