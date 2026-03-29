DROP TRIGGER IF EXISTS trg_company_updated_at ON companies;

DROP INDEX IF EXISTS idx_companies_open_vacancies_desc;
DROP INDEX IF EXISTS idx_companies_created_at_desc;

DROP TABLE IF EXISTS companies;

DROP FUNCTION IF EXISTS set_updated_at;
