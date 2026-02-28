SELECT v.id
FROM vacancies v
         JOIN companies c ON v.company_id = c.id
WHERE v.salary_from < 100000
  AND v.salary_to > 40000
  AND v.city = ANY(ARRAY['Moscow', 'London'])
  AND v.work_format = 'remote'
  AND v.duration_to_months > 3
  AND v.is_active = true
ORDER BY v.published_at DESC, v.id DESC
LIMIT 20;