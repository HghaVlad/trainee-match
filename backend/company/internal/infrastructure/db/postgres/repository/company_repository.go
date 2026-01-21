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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/update_company"
)

type CompanyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (repo *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Company, error) {
	var company entities.Company
	err := repo.db.GetContext(ctx, &company, "SELECT * FROM companies WHERE id = $1", id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", domain_errors.ErrCompanyNotFound, id)
	}

	if err != nil {
		return nil, err
	}

	return &company, err
}

func (repo *CompanyRepository) Create(ctx context.Context, company *entities.Company) error {

	_, err := repo.db.ExecContext(ctx, "INSERT INTO companies (id, name, description, website, owner_id) VALUES ($1, $2, $3, $4, $5)",
		company.ID, company.Name, company.Description, company.Website, company.OwnerID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return domain_errors.ErrCompanyAlreadyExists
		}
	}

	return err
}

// Update updates only req's non-nil fields
func (repo *CompanyRepository) Update(ctx context.Context, req *update_company.Request) error {
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

	res, err := repo.db.ExecContext(ctx, query, args...)
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
