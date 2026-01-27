package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

type CompanyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (repo *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	err := repo.db.GetContext(ctx, &company, "SELECT * FROM companies WHERE id = $1", id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", domain_errors.ErrCompanyNotFound, id)
	}

	if err != nil {
		return nil, err
	}

	return &company, err
}

func (repo *CompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	exec := repo.getExec(ctx)

	_, err := exec.ExecContext(ctx, "INSERT INTO companies (id, name, description, website, owner_id) VALUES ($1, $2, $3, $4, $5)",
		company.ID, company.Name, company.Description, company.Website, company.OwnerID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return domain_errors.ErrCompanyAlreadyExists
		}
	}

	return err
}

func (repo *CompanyRepository) ListByVacanciesCnt(
	ctx context.Context,
	cursor *list_companies.VacanciesCntCursor,
	limit int,
) (
	[]list_companies.CompanySummary,
	*list_companies.VacanciesCntCursor,
	error,
) {
	var query string
	var args []any

	if cursor == nil {
		query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		ORDER BY open_vacancies_count DESC, name
		LIMIT $1`
		args = []any{limit}
	} else {
		query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		WHERE (open_vacancies_count, name) < ($1, $2)
		ORDER BY open_vacancies_count DESC, name
		LIMIT $3`
		args = []any{cursor.Count, cursor.Name, limit}
	}

	var companies []list_companies.CompanySummary

	err := repo.db.SelectContext(ctx, &companies, query, args...)

	if err != nil {
		return nil, nil, err
	}

	if len(companies) < limit {
		return companies, nil, nil
	}

	last := companies[len(companies)-1]
	nextCursor := list_companies.VacanciesCntCursor{
		Count: last.OpenVacanciesCnt,
		Name:  last.Name,
	}

	return companies, &nextCursor, nil
}

// ListByCreatedAtDesc takes companies "after" cursor, returns them with next cursor
func (repo *CompanyRepository) ListByCreatedAtDesc(
	ctx context.Context,
	cursor *list_companies.CreatedAtCursor,
	limit int,
) (
	[]list_companies.CompanySummary,
	*list_companies.CreatedAtCursor,
	error,
) {
	var query string
	var args []any

	if cursor == nil {
		query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		ORDER BY created_at DESC, name
		LIMIT $1`
		args = []any{limit}
	} else {
		query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		WHERE (created_at, name) < ($1, $2)
		ORDER BY created_at DESC, name
		LIMIT $3`
		args = []any{cursor.CreatedAt, cursor.Name, limit}
	}

	var companies []list_companies.CompanySummary

	err := repo.db.SelectContext(ctx, &companies, query, args...)

	if err != nil {
		return nil, nil, err
	}

	if len(companies) < limit {
		return companies, nil, nil
	}

	last := companies[len(companies)-1]
	nextCursor := list_companies.CreatedAtCursor{
		CreatedAt: last.CreatedAt,
		Name:      last.Name,
	}

	return companies, &nextCursor, nil
}

func (repo *CompanyRepository) ListByName(
	ctx context.Context,
	cursor *list_companies.NameCursor,
	limit int,
) (
	[]list_companies.CompanySummary,
	*list_companies.NameCursor,
	error,
) {
	var query string
	var args []any

	if cursor == nil {
		query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		ORDER BY name
		LIMIT $1`
		args = []any{limit}
	} else {
		query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		WHERE name > $1
		ORDER BY name
		LIMIT $2`
		args = []any{cursor.Name, limit}
	}

	var companies []list_companies.CompanySummary

	err := repo.db.SelectContext(ctx, &companies, query, args...)

	if err != nil {
		return nil, nil, err
	}

	if len(companies) < limit {
		return companies, nil, nil
	}

	last := companies[len(companies)-1]
	nextCursor := list_companies.NameCursor{
		Name: last.Name,
	}

	return companies, &nextCursor, nil
}

// Update updates only req's non-nil fields
func (repo *CompanyRepository) Update(ctx context.Context, req *update_company.Request) error {
	exec := repo.getExec(ctx)
	setParts := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argID))
		args = append(args, *req.Name)
		argID++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argID))
		args = append(args, req.Description)
		argID++
	}

	if req.Website != nil {
		setParts = append(setParts, fmt.Sprintf("website = $%d", argID))
		args = append(args, req.Website)
		argID++
	}

	if len(setParts) == 0 {
		return nil // ничего не обновляем
	}

	query := fmt.Sprintf(
		"UPDATE companies SET %s WHERE id = $%d",
		strings.Join(setParts, ", "),
		argID,
	)

	args = append(args, req.ID)

	res, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain_errors.ErrCompanyAlreadyExists
		}
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain_errors.ErrCompanyNotFound
	}

	return nil
}

func (repo *CompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`DELETE FROM companies WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain_errors.ErrCompanyNotFound
	}

	return nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (repo *CompanyRepository) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value(infra_postgres.TxKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return repo.db
}
