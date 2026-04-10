ALTER TABLE vacancies
    DROP CONSTRAINT IF EXISTS vacancies_created_by_user_id_company_id_fkey;

DROP INDEX IF EXISTS idx_company_members_company;

DROP TABLE IF EXISTS company_members;

DROP TYPE IF EXISTS company_member_role_enum;
