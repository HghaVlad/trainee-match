DROP TRIGGER IF EXISTS trg_vacancy_updated_at ON vacancies;

DROP INDEX IF EXISTS idx_vacancies_feed;
DROP INDEX IF EXISTS idx_vacancies_company_feed;
DROP INDEX IF EXISTS idx_vacancies_salary_feed;
DROP INDEX IF EXISTS idx_vacancies_salary_feed_asc;

DROP TABLE IF EXISTS vacancies;

DROP TYPE IF EXISTS vacancy_status_enum;
DROP TYPE IF EXISTS employment_type_enum;
DROP TYPE IF EXISTS work_format_enum;
