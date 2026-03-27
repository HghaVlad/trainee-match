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

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
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

func (repo *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	var comp company.Company
	err := repo.db.GetContext(ctx, &comp, "SELECT * FROM companies WHERE id = $1", id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", company.ErrCompanyNotFound, id)
	}

	if err != nil {
		return nil, err
	}

	return &comp, err
}

func (repo *CompanyRepository) Create(ctx context.Context, comp *company.Company) error {
	exec := repo.getExec(ctx)

	_, err := exec.ExecContext(ctx, "INSERT INTO companies (id, name, description, website) VALUES ($1, $2, $3, $4)",
		comp.ID, comp.Name, comp.Description, comp.Website)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return company.ErrCompanyAlreadyExists
		}
	}

	return err
}

func (repo *CompanyRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM companies WHERE id = $1)`

	var exists bool
	err := repo.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *CompanyRepository) ListByVacanciesCnt(
	ctx context.Context,
	cursor *list.VacanciesCntCursor,
	limit int,
) (
	[]list.CompanySummary,
	*list.VacanciesCntCursor,
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
		WHERE open_vacancies_count < $1 OR (open_vacancies_count = $1 AND name > $2)
		ORDER BY open_vacancies_count DESC, name
		LIMIT $3`
		args = []any{cursor.Count, cursor.Name, limit}
	}

	var companies []list.CompanySummary

	err := repo.db.SelectContext(ctx, &companies, query, args...)

	if err != nil {
		return nil, nil, err
	}

	if len(companies) < limit {
		return companies, nil, nil
	}

	last := companies[len(companies)-1]
	nextCursor := list.VacanciesCntCursor{
		Count: last.OpenVacanciesCnt,
		Name:  last.Name,
	}

	return companies, &nextCursor, nil
}

// ListByCreatedAtDesc takes companies "after" cursor, returns them with next cursor
func (repo *CompanyRepository) ListByCreatedAtDesc(
	ctx context.Context,
	cursor *list.CreatedAtCursor,
	limit int,
) (
	[]list.CompanySummary,
	*list.CreatedAtCursor,
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
		WHERE created_at < $1 OR (created_at = $1 AND name > $2)
		ORDER BY created_at DESC, name
		LIMIT $3`
		args = []any{cursor.CreatedAt, cursor.Name, limit}
	}

	var companies []list.CompanySummary

	err := repo.db.SelectContext(ctx, &companies, query, args...)

	if err != nil {
		return nil, nil, err
	}

	if len(companies) < limit {
		return companies, nil, nil
	}

	last := companies[len(companies)-1]
	nextCursor := list.CreatedAtCursor{
		CreatedAt: last.CreatedAt,
		Name:      last.Name,
	}

	return companies, &nextCursor, nil
}

func (repo *CompanyRepository) ListByName(
	ctx context.Context,
	cursor *list.NameCursor,
	limit int,
) (
	[]list.CompanySummary,
	*list.NameCursor,
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

	var companies []list.CompanySummary

	err := repo.db.SelectContext(ctx, &companies, query, args...)

	if err != nil {
		return nil, nil, err
	}

	if len(companies) < limit {
		return companies, nil, nil
	}

	last := companies[len(companies)-1]
	nextCursor := list.NameCursor{
		Name: last.Name,
	}

	return companies, &nextCursor, nil
}

// Update updates only req's non-nil fields
func (repo *CompanyRepository) Update(ctx context.Context, req *update.Request) error {
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
			return company.ErrCompanyAlreadyExists
		}
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

func (repo *CompanyRepository) IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`UPDATE companies SET open_vacancies_count = open_vacancies_count+1 WHERE id = $1`, id)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

func (repo *CompanyRepository) DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`UPDATE companies SET open_vacancies_count = open_vacancies_count-1 WHERE id = $1`, id)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return company.ErrCompanyNotFound
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
		return company.ErrCompanyNotFound
	}

	return nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (repo *CompanyRepository) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return repo.db
}
