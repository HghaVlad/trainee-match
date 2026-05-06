package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (h *HrProjectionRepo) CreateIdempotent(ctx context.Context, hrProj userhr.Projection) error {
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

func (h *HrProjectionRepo) GetByUsername(ctx context.Context, username string) (*userhr.Projection, error) {
	q := postgres.GetQuerier(ctx, h.db)

	const query = `SELECT user_id, username, email, created_at
		FROM hr_user_projection
		WHERE username = $1`

	var userHrProj userhr.Projection

	err := q.QueryRow(ctx, query, username).
		Scan(&userHrProj.UserID, &userHrProj.Username, &userHrProj.Email, &userHrProj.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, userhr.ErrNotFound
		}

		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &userHrProj, nil
}
