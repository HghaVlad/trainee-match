package repository

import (
	"context"
	"encoding/json"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type AppSnapshot struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewAppSnapshot(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *AppSnapshot {
	return &AppSnapshot{
		db:     db,
		getter: getter,
	}
}

func (a *AppSnapshot) CreateIdempotent(ctx context.Context, appSnapshot application.Snapshot) error {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `INSERT INTO application_snapshots
    (id, resume_id, candidate_id, resume_name, resume_data,
     full_name, email, telegram, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT (id) DO NOTHING`

	resDataB, err := json.Marshal(appSnapshot.ResumeData)
	if err != nil {
		return fmt.Errorf("app snapshot create idempotent: marshal json: %v", err)
	}

	_, err = q.Exec(ctx, query, appSnapshot.ID, appSnapshot.ResumeID,
		appSnapshot.CandidateID, appSnapshot.ResumeName,
		resDataB, appSnapshot.FullName, appSnapshot.Email,
		appSnapshot.Telegram, appSnapshot.CreatedAt)

	if err != nil {
		return fmt.Errorf("app snapshot create idempotent: %v", err)
	}

	return nil
}
