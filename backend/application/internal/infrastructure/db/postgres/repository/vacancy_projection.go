package repository

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type VacancyProjection struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewVacancyProjection(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *VacancyProjection {
	return &VacancyProjection{
		db:     db,
		getter: getter,
	}
}

func (p *VacancyProjection) GetByID(ctx context.Context, vacID uuid.UUID) (*projection.Vacancy, error) {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		SELECT id, company_id, company_name, title, status, created_at, updated_at
		FROM vacancy_projection
		WHERE id = $1
	`

	var vacancy projection.Vacancy

	err := q.QueryRow(ctx, query, vacID).
		Scan(&vacancy.ID, &vacancy.CompanyID, &vacancy.CompanyName, &vacancy.Title,
			&vacancy.Status, &vacancy.CreatedAt, &vacancy.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, projection.ErrVacancyNotFound
		}

		return nil, fmt.Errorf("get vacancy projection: %w", err)
	}

	return &vacancy, nil
}

func (p *VacancyProjection) GetCompanyIDByVacancyID(
	ctx context.Context,
	vacancyID uuid.UUID,
) (uuid.UUID, error) {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		SELECT company_id
		FROM vacancy_projection
		WHERE id = $1
	`

	var companyID uuid.UUID

	err := q.QueryRow(ctx, query, vacancyID).
		Scan(&companyID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, projection.ErrVacancyNotFound
		}

		return uuid.Nil, fmt.Errorf("get vacancy company id: %w", err)
	}

	return companyID, nil
}

// CheckHrAccess returns companyID if success,
// projection.ErrVacancyNotFound otherwise
func (p *VacancyProjection) CheckHrAccess(
	ctx context.Context,
	userID, vacancyID uuid.UUID,
) (uuid.UUID, error) {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `SELECT v.company_id
		FROM vacancy_projection v
		JOIN company_members cm ON cm.company_id = v.company_id
		WHERE v.id = $1 AND cm.user_id = $2`

	var companyID uuid.UUID

	err := q.QueryRow(ctx, query, vacancyID, userID).Scan(&companyID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, projection.ErrVacancyNotFound
		}

		return uuid.Nil, fmt.Errorf("check hr access company id: %w", err)
	}

	return companyID, nil
}
