DROP INDEX IF EXISTS idx_status_history_app;
DROP INDEX IF EXISTS idx_applications_vacancy;
DROP INDEX IF EXISTS idx_applications_candidate;
DROP INDEX IF EXISTS idx_applications_company;
DROP INDEX IF EXISTS uniq_active_application;
DROP INDEX IF EXISTS idx_snapshots_resume;
DROP INDEX IF EXISTS idx_snapshots_candidate;
DROP INDEX IF EXISTS idx_resume_candidate;
DROP INDEX IF EXISTS idx_vacancies_company;
DROP INDEX IF EXISTS idx_company_members_user;

DROP TABLE IF EXISTS application_status_history;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS application_snapshots;
DROP TABLE IF EXISTS company_members;
DROP TABLE IF EXISTS vacancy_projection;
DROP TABLE IF EXISTS resume_projection;
DROP TABLE IF EXISTS candidate_projection;

DROP TYPE IF EXISTS resume_status_enum;
DROP TYPE IF EXISTS vacancy_status_enum;
DROP TYPE IF EXISTS application_actor_enum;
DROP TYPE IF EXISTS application_status_enum;
