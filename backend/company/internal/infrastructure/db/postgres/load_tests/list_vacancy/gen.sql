INSERT INTO companies (
    id,
    name,
    description,
    website,
    logo_key,
    open_vacancies_count
)
SELECT
    gen_random_uuid(),
    'Company ' || g,
    'Description ' || g,
    'https://company' || g || '.com',
    'logo_' || g,
    0
FROM generate_series(1, 2000) g;

WITH company_ids AS (
    SELECT array_agg(id) AS ids FROM companies
)
INSERT INTO vacancies (
    id,
    company_id,
    title,
    description,
    work_format,
    city,
    duration_from_months,
    duration_to_months,
    employment_type,
    hours_per_week_from,
    hours_per_week_to,
    flexible_schedule,
    is_paid,
    salary_from,
    salary_to,
    internship_to_offer,
    is_active,
    published_at
)
SELECT
    gen_random_uuid(),
    ids[(random() * (array_length(ids,1)-1) + 1)::int],
    'Vacancy ' || g,
    'Description',
    (ARRAY['onsite','remote','hybrid'])[floor(random()*3)+1]::work_format_enum,
    (ARRAY['Moscow','Berlin','London','NYC'])[floor(random()*4)+1],
    1,
    6,
    (ARRAY['full_time','part_time','internship'])[floor(random()*3)+1]::employment_type_enum,
    20,
    40,
    random() < 0.5,
    random() < 0.7,
    CASE WHEN random() < 0.2 THEN NULL ELSE floor(random()*200000)::int END,
    CASE WHEN random() < 0.2 THEN NULL ELSE floor(random()*300000 + 200000)::int END,
    random() < 0.3,
    true,
    now() - (random() * interval '365 days')
FROM generate_series(1, 300000) g,
     company_ids;