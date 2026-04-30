package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/getpublished"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listbycomp"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
)

type VacancyRepo struct {
	db *pgxpool.Pool
}

func NewVacancyRepo(db *pgxpool.Pool) *VacancyRepo {
	return &VacancyRepo{db: db}
}

// GetByID returns ErrVacancyNotFound if vacancy's company_id != companyID
func (repo *VacancyRepo) GetByID(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
) (*vacancy.Vacancy, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT 
    id, company_id, title, description, work_format, city,
    duration_from_days, duration_to_days, employment_type,
    hours_per_week_from, hours_per_week_to, flexible_schedule, is_paid,
    salary_from, salary_to, internship_to_offer, status, created_by_user_id,
    published_at, created_at, updated_at
	FROM vacancies 
	WHERE id = $1 AND company_id = $2`

	var vac vacancy.Vacancy

	err := q.QueryRow(ctx, query, vacancyID, companyID).
		Scan(&vac.ID, &vac.CompanyID, &vac.Title, &vac.Description, &vac.WorkFormat, &vac.City,
			&vac.DurationFromDays, &vac.DurationToDays, &vac.EmploymentType,
			&vac.HoursPerWeekFrom, &vac.HoursPerWeekTo, &vac.FlexibleSchedule, &vac.IsPaid,
			&vac.SalaryFrom, &vac.SalaryTo, &vac.InternshipToOffer, &vac.Status, &vac.CreatedBy,
			&vac.PublishedAt, &vac.CreatedAt, &vac.UpdatedAt,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", vacancy.ErrVacancyNotFound, vacancyID)
	}

	if err != nil {
		return nil, fmt.Errorf("get vacancy: %w", err)
	}

	return &vac, nil
}

func (repo *VacancyRepo) GetPublishedByID(ctx context.Context, vacancyID uuid.UUID) (*getpublished.Response, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT
    v.id, v.company_id, c.name,
	v.title, v.description, v.work_format, v.city,
	v.duration_from_days, v.duration_to_days,
	v.employment_type, v.hours_per_week_from, v.hours_per_week_to,
	v.flexible_schedule, v.is_paid, v.salary_from, v.salary_to,
	v.internship_to_offer, v.published_at
	FROM vacancies v
	JOIN companies c ON c.id = v.company_id
	WHERE v.id = $1 AND v.status = 'published'`

	var vac getpublished.Response

	err := q.QueryRow(ctx, query, vacancyID).
		Scan(
			&vac.ID, &vac.CompanyID, &vac.CompanyName,
			&vac.Title, &vac.Description, &vac.WorkFormat, &vac.City,
			&vac.DurationFromDays, &vac.DurationToDays,
			&vac.EmploymentType, &vac.HoursPerWeekFrom, &vac.HoursPerWeekTo,
			&vac.FlexibleSchedule, &vac.IsPaid, &vac.SalaryFrom, &vac.SalaryTo,
			&vac.InternshipToOffer, &vac.PublishedAt,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", vacancy.ErrVacancyNotFound, vacancyID)
	}

	if err != nil {
		return nil, fmt.Errorf("get published vacancy: %w", err)
	}

	return &vac, nil
}

func (repo *VacancyRepo) GetPublishedEventView(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
) (*publish.PublishedEventView, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT 
    v.id, v.company_id, v.title, v.status, c.name
	FROM vacancies v
	JOIN companies c ON v.company_id = c.id
	WHERE v.id = $1 AND v.company_id = $2`

	var vac publish.PublishedEventView

	err := q.QueryRow(ctx, query, vacancyID, companyID).Scan(
		&vac.ID, &vac.CompanyID, &vac.Title, &vac.Status, &vac.CompanyName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get published event view: %w", vacancy.ErrVacancyNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("get published event view: %w", err)
	}

	return &vac, nil
}

func (repo *VacancyRepo) Create(ctx context.Context, vacancy *vacancy.Vacancy) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `INSERT INTO vacancies (
			id, company_id, title, description,	work_format, city,
			duration_from_days, duration_to_days,
		    employment_type, hours_per_week_from, hours_per_week_to,
			flexible_schedule, is_paid, salary_from, salary_to,
			internship_to_offer, status, created_by_user_id, published_at
		) VALUES (
		    $1, $2, $3, $4,	$5, $6, 
			$7, $8,
			$9,	$10, $11,
			$12, $13, $14, $15,
			$16, $17, $18, $19
		)`

	_, err := q.Exec(ctx, query,
		vacancy.ID, vacancy.CompanyID, vacancy.Title, vacancy.Description, vacancy.WorkFormat, vacancy.City,
		vacancy.DurationFromDays, vacancy.DurationToDays,
		vacancy.EmploymentType, vacancy.HoursPerWeekFrom, vacancy.HoursPerWeekTo,
		vacancy.FlexibleSchedule, vacancy.IsPaid, vacancy.SalaryFrom, vacancy.SalaryTo,
		vacancy.InternshipToOffer, vacancy.Status, vacancy.CreatedBy, vacancy.PublishedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return company.ErrCompanyNotFound
		}
	}

	if err != nil {
		return fmt.Errorf("create vacancy: %w", err)
	}

	return nil
}

// ListPublishedSummaries uses dynamic sql for all the requirements (filters).
// Pass cursor as a pointer
func (repo *VacancyRepo) ListPublishedSummaries(
	ctx context.Context,
	requirements *list.Requirements,
	order list.Order,
	cursor any,
	limit int,
) (
	[]list.VacancySummary, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	filtersCondition, args := listVacRequirementsToSQL(requirements)
	if filtersCondition != "" {
		filtersCondition = "AND " + filtersCondition
	}

	cursorCondition := ""
	if cursor != nil && !reflect.ValueOf(cursor).IsNil() {
		cursorCondition, args = listVacCursorToSQL(order, cursor, args)
		cursorCondition = "AND " + cursorCondition
	}

	if order == list.OrderSalaryDesc || order == list.OrderSalaryAsc {
		filtersCondition += andSalaryNotNull
	}

	orderBy := listVacOrderToSQL(order)

	args = append(args, limit)

	const query = `SELECT 
    v.id, v.company_id, c.name, v.title, v.work_format,
    v.city, v.employment_type, v.is_paid,
    v.salary_from, v.salary_to, v.published_at
	FROM vacancies v
	JOIN companies c ON v.company_id = c.id
	WHERE v.status = 'published' %s %s
	ORDER BY %s
	LIMIT $%d`

	filledQuery := fmt.Sprintf(query, filtersCondition, cursorCondition, orderBy, len(args))

	rows, err := q.Query(ctx, filledQuery, args...)

	if err != nil {
		return nil, fmt.Errorf("list published vacancy summaries: %w", err)
	}

	var vacancies []list.VacancySummary

	for rows.Next() {
		var vac list.VacancySummary

		err := rows.Scan(
			&vac.ID, &vac.CompanyID, &vac.CompanyName, &vac.Title, &vac.WorkFormat,
			&vac.City, &vac.EmploymentType, &vac.IsPaid,
			&vac.SalaryFrom, &vac.SalaryTo, &vac.PublishedAt)

		if err != nil {
			return nil, fmt.Errorf("list published vacancy summaries scan: %w", err)
		}

		vacancies = append(vacancies, vac)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list published vacancy summaries rows error: %w", err)
	}

	return vacancies, nil
}

func (repo *VacancyRepo) ListByCompanySummaries(
	ctx context.Context,
	compID uuid.UUID,
	requirements *list.Requirements,
	status *vacancy.Status,
	cursor *listbycomp.CreatedAtCursor,
	limit int,
) ([]listbycomp.VacancySummary, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	filtersCondition, args := listVacRequirementsToSQL(requirements)
	if filtersCondition != "" {
		filtersCondition = "AND " + filtersCondition
	}

	statusCondition, args := listByCompStatusToSQL(status, args)
	if statusCondition != "" {
		statusCondition = "AND " + statusCondition
	}

	cursorCondition := ""
	if cursor != nil {
		cursorCondition, args = listByCompCreatedAtCursorToSQL(*cursor, args)
		cursorCondition = "AND " + cursorCondition
	}

	args = append(args, compID, limit)

	const query = `SELECT
    v.id, v.title, v.work_format,
    v.city, v.employment_type, v.is_paid,
    v.salary_from, v.salary_to, v.status, v.created_at
	FROM vacancies v
	WHERE 1=1 %s %s %s
		AND v.company_id = $%d
	ORDER BY v.created_at DESC, v.id DESC
	LIMIT $%d`

	filledQuery := fmt.Sprintf(query, filtersCondition, statusCondition, cursorCondition, len(args)-1, len(args))

	rows, err := q.Query(ctx, filledQuery, args...)

	if err != nil {
		return nil, fmt.Errorf("list vacancy by company: %w", err)
	}

	var vacancies []listbycomp.VacancySummary

	for rows.Next() {
		var vac listbycomp.VacancySummary

		err := rows.Scan(
			&vac.ID, &vac.Title, &vac.WorkFormat,
			&vac.City, &vac.EmploymentType, &vac.IsPaid,
			&vac.SalaryFrom, &vac.SalaryTo, &vac.Status, &vac.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("list vacancy by company: %w", err)
		}

		vacancies = append(vacancies, vac)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list vacancy by company rows error: %w", err)
	}

	return vacancies, nil
}

func (repo *VacancyRepo) Update(ctx context.Context, v *vacancy.Vacancy) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `UPDATE vacancies SET
			title = $1,	description = $2, work_format = $3,	city = $4,
			duration_from_days = $5, duration_to_days = $6,
			employment_type = $7, hours_per_week_from = $8,
			hours_per_week_to = $9,	flexible_schedule = $10,
			is_paid = $11, salary_from = $12, salary_to = $13,
			internship_to_offer = $14, updated_at = now()
		WHERE id = $15`

	cmd, err := q.Exec(ctx, query,
		v.Title, v.Description, v.WorkFormat, v.City,
		v.DurationFromDays, v.DurationToDays,
		v.EmploymentType, v.HoursPerWeekFrom,
		v.HoursPerWeekTo, v.FlexibleSchedule,
		v.IsPaid, v.SalaryFrom, v.SalaryTo,
		v.InternshipToOffer,
		v.ID,
	)

	if err != nil {
		return fmt.Errorf("update vacancy: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

func (repo *VacancyRepo) Publish(ctx context.Context, vacID, compID uuid.UUID) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `UPDATE vacancies
		SET status = 'published', published_at = now()
		WHERE id = $1 AND company_id = $2`

	cmd, err := q.Exec(ctx, query, vacID, compID)

	if err != nil {
		return fmt.Errorf("publish vacancy: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

func (repo *VacancyRepo) Archive(ctx context.Context, vacID, compID uuid.UUID) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `UPDATE vacancies
		SET status = 'archived', published_at = NULL
		WHERE id = $1 AND company_id = $2`

	cmd, err := q.Exec(ctx, query, vacID, compID)

	if err != nil {
		return fmt.Errorf("archive vacancy: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

func (repo *VacancyRepo) Delete(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `DELETE FROM vacancies WHERE id = $1 AND company_id = $2`

	cmd, err := q.Exec(ctx, query, vacancyID, companyID)

	if err != nil {
		return fmt.Errorf("delete vacancy: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}
