package repository

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type AppStatusHistoryRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewAppStatusHistoryRepo(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *AppStatusHistoryRepo {
	return &AppStatusHistoryRepo{
		db:     db,
		getter: getter,
	}
}

func (r *AppStatusHistoryRepo) Add(ctx context.Context, change application.StatusChange) error {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `INSERT INTO application_status_history
		(id, application_id, status,
		changed_by_user_id, changed_by_role, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := q.Exec(ctx, query, change.ID, change.ApplicationID, change.Status,
		change.ChangedByUserID, change.ChangedByRole, change.Comment, change.CreatedAt)

	if err != nil {
		return fmt.Errorf("create application status history: %w", err)
	}

	return nil
}
