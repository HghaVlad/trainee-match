package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
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
