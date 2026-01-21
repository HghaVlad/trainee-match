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
