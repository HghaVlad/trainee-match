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

func (v *VacancyProjection) GetByID(ctx context.Context, vacID uuid.UUID) (*projection.Vacancy, error) {
	q := v.getter.DefaultTrOrDB(ctx, v.db)

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
