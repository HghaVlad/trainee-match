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

type CandidateProjection struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewCandidateProjection(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *CandidateProjection {
	return &CandidateProjection{
		db:     db,
		getter: getter,
	}
}

func (c *CandidateProjection) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*projection.Candidate, error) {
	q := c.getter.DefaultTrOrDB(ctx, c.db)

	const query = `SELECT id, full_name, email, telegram, created_at, updated_at
		FROM candidate_projection
		WHERE id = $1`

	var candProj projection.Candidate

	err := q.QueryRow(ctx, query, userID).
		Scan(&candProj.ID, &candProj.FullName, &candProj.Email, &candProj.Telegram,
			&candProj.CreatedAt, &candProj.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, projection.ErrCandidateNotFound
		}

		return nil, fmt.Errorf("get candidate projection: %v", err)
	}

	return &candProj, nil
}
