package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type ResumeProjection struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewResumeProjection(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *ResumeProjection {
	return &ResumeProjection{
		db:     db,
		getter: getter,
	}
}

func (r *ResumeProjection) GetByID(ctx context.Context, id uuid.UUID) (*projection.Resume, error) {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT id, candidate_id, name, data, status, created_at, updated_at
		FROM resume_projection
		WHERE id = $1
	`

	var resume projection.Resume
	var dataRaw []byte

	err := q.QueryRow(ctx, query, id).
		Scan(&resume.ID, &resume.CandidateID, &resume.Name, &dataRaw, &resume.Status,
			&resume.CreatedAt, &resume.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, projection.ErrResumeNotFound
		}

		return nil, fmt.Errorf("get resume projection: %w", err)
	}

	if err := json.Unmarshal(dataRaw, &resume.Data); err != nil {
		return nil, fmt.Errorf("unmarshal resume data: %w", err)
	}

	return &resume, nil
}
