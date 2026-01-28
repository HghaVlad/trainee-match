package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
)

type VacancyRepo struct {
	db *sqlx.DB
}

func NewVacancyRepo(db *sqlx.DB) *VacancyRepo {
	return &VacancyRepo{db: db}
}

// GetByID returns ErrVacancyNotFound if vacancy's company_id != companyID
func (repo *VacancyRepo) GetByID(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error) {
	var company domain.Vacancy
	err := repo.db.GetContext(ctx, &company, "SELECT * FROM vacancies WHERE id = $1 AND company_id = $2",
		vacancyID, companyID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", domain_errors.ErrVacancyNotFound, vacancyID)
	}

	if err != nil {
		return nil, err
	}

	return &company, err
}

func (repo *VacancyRepo) Create(ctx context.Context, vacancy *domain.Vacancy) error {
	exec := repo.getExec(ctx)

	_, err := exec.ExecContext(ctx, `
		INSERT INTO vacancies (
			id, company_id,	title, description,	work_format, city,
			duration_from_months, duration_to_months,
		    employment_type, hours_per_week_from, hours_per_week_to,
			flexible_schedule, is_paid, salary_from, salary_to,
			internship_to_offer
		) VALUES (
		    $1, $2, $3, $4,	$5, $6,
			$7, $8,
			$9,	$10, $11,
			$12, $13, $14, $15,
			$16
		)
	`,
		vacancy.ID, vacancy.CompanyID, vacancy.Title, vacancy.Description, vacancy.WorkFormat, vacancy.City,
		vacancy.DurationFromMonths, vacancy.DurationToMonths,
		vacancy.EmploymentType, vacancy.HoursPerWeekFrom, vacancy.HoursPerWeekTo,
		vacancy.FlexibleSchedule, vacancy.IsPaid, vacancy.SalaryFrom, vacancy.SalaryTo,
		vacancy.InternshipToOffer,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return domain_errors.ErrCompanyNotFound
		}
	}

	return err
}

// returns sqlx.TX if we're in transaction or r.db if not
func (repo *VacancyRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value(infra_postgres.TxKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return repo.db
}
