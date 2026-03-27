package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/getpublished"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listbycomp"
)

type VacancyRepo struct {
	db *sqlx.DB
}

func NewVacancyRepo(db *sqlx.DB) *VacancyRepo {
	return &VacancyRepo{db: db}
}

// GetByID returns ErrVacancyNotFound if vacancy's company_id != companyID
func (repo *VacancyRepo) GetByID(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
) (*vacancy.Vacancy, error) {

	var vac vacancy.Vacancy
	err := repo.db.GetContext(ctx, &vac,
		"SELECT * FROM vacancies WHERE id = $1 AND company_id = $2",
		vacancyID, companyID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", vacancy.ErrVacancyNotFound, vacancyID)
	}

	if err != nil {
		return nil, err
	}

	return &vac, err
}

func (repo *VacancyRepo) GetPublishedByID(ctx context.Context, vacancyID uuid.UUID) (*getpublished.Response, error) {
	var vac getpublished.Response
	err := repo.db.GetContext(ctx, &vac, `
		SELECT v.id, v.company_id, c.name AS company_name,
		       v.title, v.description, v.work_format, v.city,
		       v.duration_from_days, v.duration_to_days,
		       v.employment_type, v.hours_per_week_from, v.hours_per_week_to,
		       v.flexible_schedule, v.is_paid, v.salary_from, v.salary_to,
		       v.internship_to_offer, v.published_at
		FROM vacancies v
		JOIN companies c ON c.id = v.company_id
		WHERE v.id = $1 AND v.status = 'published'
	`, vacancyID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", vacancy.ErrVacancyNotFound, vacancyID)
	}

	if err != nil {
		return nil, err
	}

	return &vac, nil
}

func (repo *VacancyRepo) Create(ctx context.Context, vacancy *vacancy.Vacancy) error {
	exec := repo.getExec(ctx)

	_, err := exec.ExecContext(ctx, `
		INSERT INTO vacancies (
			id, company_id, created_by_user_id,	title, description,	work_format, city,
			duration_from_days, duration_to_days,
		    employment_type, hours_per_week_from, hours_per_week_to,
			flexible_schedule, is_paid, salary_from, salary_to,
			internship_to_offer, status, published_at
		) VALUES (
		    $1, $2, $3, $4,	$5, $6, 
			$7, $8,
			$9,	$10, $11,
			$12, $13, $14, $15,
			$16, $17, $18, $19
		)
	`,
		vacancy.ID, vacancy.CompanyID, vacancy.CreatedBy, vacancy.Title, vacancy.Description, vacancy.WorkFormat, vacancy.City,
		vacancy.DurationFromDays, vacancy.DurationToDays,
		vacancy.EmploymentType, vacancy.HoursPerWeekFrom, vacancy.HoursPerWeekTo,
		vacancy.FlexibleSchedule, vacancy.IsPaid, vacancy.SalaryFrom, vacancy.SalaryTo,
		vacancy.InternshipToOffer, vacancy.Status, vacancy.PublishedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return company.ErrCompanyNotFound
		}
	}

	return err
}

// ListPublished uses dynamic sql for all the requirements (filters).
// Pass cursor as a pointer
func (repo *VacancyRepo) ListPublished(
	ctx context.Context,
	requirements *list.Requirements,
	order list.Order,
	cursor any,
	limit int,
) (
	[]list.VacancySummary, error) {

	requireFilters, args := listVacRequirementsToSQL(requirements)

	cursorCondition := ""
	if cursor != nil && !reflect.ValueOf(cursor).IsNil() {
		cursorCondition, args = listVacCursorToSQL(order, cursor, args)
		cursorCondition = "AND " + cursorCondition
	}

	if order == list.OrderSalaryDesc || order == list.OrderSalaryAsc {
		requireFilters += andSalaryNotNull
	}

	orderBy := listVacOrderToSQL(order)

	args = append(args, limit)

	query := fmt.Sprintf(`SELECT v.id, v.company_id, c.name AS company_name, v.title, v.work_format, v.city, v.employment_type,
       	v.is_paid, v.salary_from, v.salary_to, v.published_at
		FROM vacancies v
		JOIN companies c ON v.company_id = c.id
		WHERE %s %s AND v.status = 'published'
		%s 
		LIMIT $%d`, requireFilters, cursorCondition, orderBy, len(args))

	var vacancies []list.VacancySummary

	err := repo.db.SelectContext(ctx, &vacancies, query, args...)
	return vacancies, err
}

func (repo *VacancyRepo) ListByCompanyByPublishedAt(
	ctx context.Context,
	compID uuid.UUID,
	cursor *listbycomp.PublishedAtCursor,
	limit int,
) (
	[]listbycomp.VacancySummary, error) {

	var query string
	var args []any

	if cursor == nil {
		query =
			`SELECT v.id, v.title, v.work_format, v.city, v.employment_type,
       	v.is_paid, v.salary_from, v.salary_to, v.published_at
		FROM vacancies v
		WHERE v.company_id = $1 AND v.status = 'published'
		ORDER BY v.published_at DESC, v.id
		LIMIT $2`
		args = []any{compID, limit}
	} else {
		query =
			`SELECT v.id, v.title, v.work_format, v.city, v.employment_type,
       	v.is_paid, v.salary_from, v.salary_to, v.published_at
		FROM vacancies v
		JOIN companies c ON v.company_id = c.id 
		WHERE v.company_id = $1 AND 
		      (v.published_at < $2 OR (v.published_at = $2 AND v.id < $3))
		  		AND v.status = 'published'
		ORDER BY v.published_at DESC, v.id DESC
		LIMIT $4`
		args = []any{compID, cursor.PublishedAt, cursor.Id, limit}
	}

	var vacancies []listbycomp.VacancySummary

	err := repo.db.SelectContext(ctx, &vacancies, query, args...)
	return vacancies, err
}

func (repo *VacancyRepo) Update(ctx context.Context, v *vacancy.Vacancy) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`UPDATE vacancies SET
			title = $1,
			description = $2,

			work_format = $3,
			city = $4,

			duration_from_days = $5,
			duration_to_days = $6,

			employment_type = $7,
			hours_per_week_from = $8,
			hours_per_week_to = $9,

			flexible_schedule = $10,

			is_paid = $11,
			salary_from = $12,
			salary_to = $13,

			internship_to_offer = $14,

			updated_at = now()
		WHERE id = $15
	`,
		v.Title,
		v.Description,
		v.WorkFormat,
		v.City,
		v.DurationFromDays,
		v.DurationToDays,
		v.EmploymentType,
		v.HoursPerWeekFrom,
		v.HoursPerWeekTo,
		v.FlexibleSchedule,
		v.IsPaid,
		v.SalaryFrom,
		v.SalaryTo,
		v.InternshipToOffer,
		v.ID,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

func (repo *VacancyRepo) Publish(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`UPDATE vacancies
		SET status = 'published', published_at = now()
		WHERE id = $1 AND company_id = $2`,
		vacID, compID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

func (repo *VacancyRepo) Archive(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`UPDATE vacancies
		SET status = 'archived', published_at = NULL
		WHERE id = $1 AND company_id = $2`,
		vacID, compID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

func (repo *VacancyRepo) Delete(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`DELETE FROM vacancies WHERE id = $1 AND company_id = $2`, vacancyID, companyID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return vacancy.ErrVacancyNotFound
	}

	return nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (repo *VacancyRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return repo.db
}
