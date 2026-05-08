package repository

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
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

func (r *AppStatusHistoryRepo) GetHistoryCandiView(
	ctx context.Context,
	appID, candID uuid.UUID,
) ([]views.StatusChangeCandidateFullView, error) {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `SELECT sh.status, sh.changed_by_role, sh.created_at, sh.comment
			FROM application_status_history sh
			JOIN applications a ON sh.application_id = a.id
			WHERE application_id = $1 AND a.candidate_id = $2`

	rows, err := q.Query(ctx, query, appID, candID)

	if err != nil {
		return nil, fmt.Errorf("get app status history candidate view: %w", err)
	}

	defer rows.Close()
	var history []views.StatusChangeCandidateFullView

	for rows.Next() {
		var statusChange views.StatusChangeCandidateFullView

		err := rows.Scan(&statusChange.Status, &statusChange.ChangedByRole,
			&statusChange.CreatedAt, &statusChange.Comment)
		if err != nil {
			return nil, fmt.Errorf("get app status history candidate view: %w", err)
		}

		history = append(history, statusChange)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get app status history candidate view: %w", err)
	}

	return history, nil
}

func (r *AppStatusHistoryRepo) GetHistoryHrView(
	ctx context.Context,
	appID, hrID uuid.UUID,
) ([]views.StatusChangeHrFullView, error) {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		SELECT sh.status, sh.created_at, sh.comment, sh.changed_by_role, sh.changed_by_user_id
		FROM application_status_history sh
		JOIN applications a ON sh.application_id = a.id
		JOIN company_members cm ON cm.company_id = a.company_id
		WHERE application_id = $1 AND cm.user_id = $2`

	rows, err := q.Query(ctx, query, appID, hrID)

	if err != nil {
		return nil, fmt.Errorf("get app status history hr view: %w", err)
	}

	defer rows.Close()
	var history []views.StatusChangeHrFullView

	for rows.Next() {
		var statusChange views.StatusChangeHrFullView

		err := rows.Scan(&statusChange.Status, &statusChange.CreatedAt, &statusChange.Comment,
			&statusChange.ChangedByRole, &statusChange.ChangedByUserID)
		if err != nil {
			return nil, fmt.Errorf("get app status history hr view: %w", err)
		}

		history = append(history, statusChange)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get app status history hr view: %w", err)
	}

	return history, nil
}
