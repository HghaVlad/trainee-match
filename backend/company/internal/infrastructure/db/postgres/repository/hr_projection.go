package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

type HrProjectionRepo struct {
	db *pgxpool.Pool
}

func NewHrProjectionRepo(db *pgxpool.Pool) *HrProjectionRepo {
	return &HrProjectionRepo{db: db}
}

func (h HrProjectionRepo) CreateIdempotent(ctx context.Context, hrProj userhr.Projection) error {
	q := postgres.GetQuerier(ctx, h.db)

	const query = `INSERT INTO hr_user_projection
    	(user_id, username, email, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING`

	_, err := q.Exec(ctx, query, hrProj.UserID, hrProj.Username, hrProj.Email, hrProj.CreatedAt)
	if err != nil {
		return fmt.Errorf("create idempotent: %w", err)
	}

	return nil
}
