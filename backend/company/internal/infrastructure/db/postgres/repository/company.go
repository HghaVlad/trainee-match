package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (repo *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT id, name, description, website,
    	logo_key, open_vacancies_count, created_at, updated_at
		FROM companies WHERE id = $1`

	var comp company.Company
	err := q.QueryRow(ctx, query, id).
		Scan(&comp.ID, &comp.Name, &comp.Description,
			&comp.Website, &comp.LogoKey, &comp.OpenVacanciesCnt,
			&comp.CreatedAt, &comp.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", company.ErrCompanyNotFound, id)
	}

	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	return &comp, nil
}

func (repo *CompanyRepository) Create(ctx context.Context, comp *company.Company) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = "INSERT INTO companies (id, name, description, website) VALUES ($1, $2, $3, $4)"

	_, err := q.Exec(ctx, query, comp.ID, comp.Name, comp.Description, comp.Website)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return company.ErrCompanyAlreadyExists
		}
	}

	if err != nil {
		return fmt.Errorf("create company: %w", err)
	}

	return nil
}

// ListSummaries lists company summaries by order after the cursor using dynamic SQL.
// Pass cursor as a pointer
func (repo *CompanyRepository) ListSummaries(
	ctx context.Context,
	order list.Order,
	cursor any,
	limit int,
) ([]list.CompanySummary, error) {
	args := make([]any, 0)

	cursorCondition := "1=1"
	if cursor != nil && !reflect.ValueOf(cursor).IsNil() {
		cursorCondition, args = listCompCursorToSQL(cursor, args)
	}

	orderBy := listCompanySummariesOrderToSQL(order)

	args = append(args, limit)

	const query = `SELECT id, name, open_vacancies_count, logo_key, created_at
		FROM companies
		WHERE %s
		ORDER BY %s
		LIMIT $%d`

	filledQuery := fmt.Sprintf(query, cursorCondition, orderBy, len(args))

	q := postgres.GetQuerier(ctx, repo.db)

	rows, err := q.Query(ctx, filledQuery, args...)

	if err != nil {
		return nil, fmt.Errorf("list company summary: %w", err)
	}

	defer rows.Close()

	var companies []list.CompanySummary

	for rows.Next() {
		var comp list.CompanySummary

		err := rows.Scan(&comp.ID, &comp.Name, &comp.OpenVacanciesCnt, &comp.LogoKey, &comp.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("list company summary scan: %w", err)
		}

		companies = append(companies, comp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return companies, nil
}

func (repo *CompanyRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT EXISTS (SELECT 1 FROM companies WHERE id = $1)`

	var exists bool
	err := q.QueryRow(ctx, query, id).Scan(&exists)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("company exists: %w", err)
	}

	return exists, nil
}

// Update updates only req's non-nil fields
func (repo *CompanyRepository) Update(ctx context.Context, req *update.Request) error {
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
		args = append(args, *req.Description)
		argID++
	}

	if req.Website != nil {
		setParts = append(setParts, fmt.Sprintf("website = $%d", argID))
		args = append(args, *req.Website)
		argID++
	}

	if len(setParts) == 0 {
		return nil // no update
	}

	query := fmt.Sprintf(
		"UPDATE companies SET %s WHERE id = $%d",
		strings.Join(setParts, ", "),
		argID,
	)

	args = append(args, req.ID)

	q := postgres.GetQuerier(ctx, repo.db)

	cmd, err := q.Exec(ctx, query, args...)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return company.ErrCompanyAlreadyExists
		}
		return fmt.Errorf("update company: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

func (repo *CompanyRepository) IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = "UPDATE companies SET open_vacancies_count = open_vacancies_count+1 WHERE id = $1"

	cmd, err := q.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("increment open_vacancies_count: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

func (repo *CompanyRepository) DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = "UPDATE companies SET open_vacancies_count = GREATEST(open_vacancies_count - 1, 0) WHERE id = $1"

	cmd, err := q.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("decrement open_vacancies_count: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

func (repo *CompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = "DELETE FROM companies WHERE id = $1"

	cmd, err := q.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete company: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

// returns sql condition, updated slice of args
func listCompCursorToSQL(cursor any, args []any) (string, []any) {
	switch c := cursor.(type) {
	case *list.VacanciesCntCursor:
		return compVacanciesCntCursorToSQL(*c, args)
	case *list.CreatedAtCursor:
		return compCreatedAtCursorToSQL(*c, args)
	case *list.NameCursor:
		return compNameCursorToSQL(*c, args)
	}

	return "", args
}

func compVacanciesCntCursorToSQL(cursor list.VacanciesCntCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"open_vacancies_count < $%d OR (open_vacancies_count = $%d AND name > $%d)",
		len(args)+1, len(args)+1, len(args)+2)

	args = append(args, cursor.Count, cursor.Name)
	return condition, args
}

func compCreatedAtCursorToSQL(cursor list.CreatedAtCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"created_at < $%d OR (created_at = $%d AND name > $%d)",
		len(args)+1, len(args)+1, len(args)+2)

	args = append(args, cursor.CreatedAt, cursor.Name)
	return condition, args
}

func compNameCursorToSQL(cursor list.NameCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"name > $%d",
		len(args)+1)

	args = append(args, cursor.Name)
	return condition, args
}

func listCompanySummariesOrderToSQL(order list.Order) string {
	switch order {
	case list.OrderVacanciesDesc:
		return "open_vacancies_count DESC, name"
	case list.OrderCreatedAtDesc:
		return "created_at DESC, name"
	case list.OrderNameAsc:
		return "name"
	}

	return ""
}
