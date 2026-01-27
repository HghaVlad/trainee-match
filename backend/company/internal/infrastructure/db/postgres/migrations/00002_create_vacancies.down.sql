DROP TRIGGER IF EXISTS trg_vacancies_updated_at ON vacancies;

DROP TABLE IF EXISTS vacancies;

DROP TYPE IF EXISTS employment_type_enum;
DROP TYPE IF EXISTS work_format_enum;