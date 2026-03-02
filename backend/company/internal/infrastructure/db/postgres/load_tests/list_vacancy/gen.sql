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

INSERT INTO company_members (
    user_id,
    company_id,
    role
)
SELECT
    gen_random_uuid(),
    id,
    'admin'::company_member_role_enum
FROM companies;



WITH expanded AS (
    SELECT
        cm.company_id,
        cm.user_id,
        gs,

        -- title
        (ARRAY[
            'Backend Developer','Frontend Developer','QA Engineer','DevOps Engineer',
            'Product Manager','Data Analyst','Mobile Developer','UI/UX Designer',
            'System Engineer','ML Engineer'
            ])[floor(random()*10 + 1)] AS base_title,

        -- work format
        (ARRAY['onsite','remote','hybrid']::work_format_enum[])
            [floor(random()*3 + 1)] AS work_format,

        -- city
        (ARRAY['Moscow','Berlin','London','New York',
            'Warsaw','Prague','Amsterdam','Dubai'])
            [floor(random()*8 + 1)] AS city,

        -- duration
        (30 + floor(random()*120))::int AS duration_from,

        -- employment
        (ARRAY['full_time','part_time','internship']::employment_type_enum[])
            [floor(random()*3 + 1)] AS employment_type,

        random() AS r_paid,
        random() AS r_salary,

        random() AS r_hours,

        random() < 0.3 AS flexible_schedule,
        random() < 0.2 AS internship_to_offer

    FROM company_members cm
             CROSS JOIN generate_series(1,150) gs
),

     prepared AS (
         SELECT
             e.*,

             CASE
                 WHEN employment_type = 'full_time' THEN 35 + floor(random()*6)::int   -- 35–40
                 WHEN employment_type = 'part_time' THEN 15 + floor(random()*11)::int  -- 15–25
                 ELSE 20 + floor(random()*11)::int                                      -- internship 20–30
                 END AS hours_from,

             CASE
                 WHEN employment_type = 'full_time' THEN 40
                 WHEN employment_type = 'part_time' THEN 30
                 ELSE 30
                 END AS hours_to,

             (e.r_paid > 0.2) AS paid_flag,

             -- базовая зарплата
             CASE
                 WHEN e.r_salary < 0.35 THEN 60000 + floor(random()*30000)
                 WHEN e.r_salary < 0.65 THEN 90000 + floor(random()*40000)
                 WHEN e.r_salary < 0.90 THEN 130000 + floor(random()*70000)
                 ELSE 200000 + floor(random()*150000)
                 END AS base_salary

         FROM expanded e
     ),

     numbered AS (
         SELECT *,
                ROW_NUMBER() OVER (ORDER BY company_id, gs) AS rn
         FROM prepared
     )

INSERT INTO vacancies (
    id,
    company_id,
    created_by_user_id,
    title,
    description,
    work_format,
    city,
    duration_from_days,
    duration_to_days,
    employment_type,
    hours_per_week_from,
    hours_per_week_to,
    flexible_schedule,
    is_paid,
    salary_from,
    salary_to,
    internship_to_offer,
    status,
    published_at
)
SELECT
    gen_random_uuid(),
    n.company_id,
    n.user_id,
    n.base_title || ' #' || n.rn,
    'We are looking for a talented specialist to join our team.',
    n.work_format,
    n.city,
    n.duration_from,
    n.duration_from + floor(random()*60)::int,
    n.employment_type,
    n.hours_from,
    n.hours_to,
    n.flexible_schedule,
    n.paid_flag,
    CASE WHEN n.paid_flag THEN n.base_salary ELSE NULL END,
    CASE WHEN n.paid_flag
             THEN n.base_salary + floor(n.base_salary * (0.10 + random()*0.20))::int
         ELSE NULL
        END,
    n.internship_to_offer,
    'published'::vacancy_status_enum,
    now() - (random() * interval '365 days')
FROM numbered n;

