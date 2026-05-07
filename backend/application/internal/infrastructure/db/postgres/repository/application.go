package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

type ApplicationRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewApplication(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *ApplicationRepo {
	return &ApplicationRepo{
		db:     db,
		getter: getter,
	}
}

func (a *ApplicationRepo) Create(ctx context.Context, app application.Application) error {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `INSERT INTO applications
	(id, resume_id, candidate_id, vacancy_id, company_id,
	 snapshot_id, status, cover_letter, created_at, updated_at)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := q.Exec(ctx, query, app.ID, app.ResumeID, app.CandidateID, app.VacancyID, app.CompanyID,
		app.SnapshotID, app.Status, app.CoverLetter, app.CreatedAt, app.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "uniq_active_application" {
			return application.ErrActiveAlreadyExists
		}
	}

	return nil
}

func (a *ApplicationRepo) GetByIDCandidateViewWithDetails(
	ctx context.Context,
	appID, candID uuid.UUID,
) (*views.CandidateViewWithDetails, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	// TODO: hat to do when vacancy is gone?

	const query = `SELECT a.id, a.status, a.cover_letter, a.created_at, a.updated_at,
        		v.id, v.company_id, v.company_name, v.title,
        		s.resume_data, s.email, s.full_name, s.telegram, s.created_at,
        		COALESCE(
					json_agg(json_build_object(
							'status', h.status,
							'changed_by_role', h.changed_by_role,
							'created_at', h.created_at
						)
						ORDER BY h.created_at
					) FILTER (WHERE h.id IS NOT NULL), '[]'
				) AS status_history
		FROM applications a
		LEFT JOIN vacancy_projection v ON v.id = a.vacancy_id
		JOIN application_snapshots s ON a.snapshot_id = s.id
		JOIN application_status_history h ON h.application_id = a.id
 		WHERE a.id = $1 AND a.candidate_id = $2
 		GROUP BY a.id, a.status, a.cover_letter, a.created_at, a.updated_at,
        		v.id, v.company_id, v.company_name, v.title,
        		s.resume_data, s.email, s.full_name, s.telegram, s.created_at`

	var (
		view             views.CandidateViewWithDetails
		resumeDataRaw    []byte
		statusHistoryRaw []byte
	)

	err := q.QueryRow(ctx, query, appID, candID).
		Scan(&view.AppID, &view.Status, &view.CoverLetter, &view.CreatedAt, &view.UpdatedAt,
			&view.VacancyID, &view.CompanyID, &view.CompanyName, &view.VacancyTitle,
			&resumeDataRaw, &view.Snapshot.Email, &view.Snapshot.FullName, &view.Snapshot.Telegram,
			&view.Snapshot.CreatedAt, &statusHistoryRaw)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrNotFound
		}

		return nil, fmt.Errorf(
			"get application candidate view with details: %w",
			err,
		)
	}

	err = json.Unmarshal(resumeDataRaw, &view.Snapshot.ResumeData)
	if err != nil {
		return nil, fmt.Errorf(
			"unmarshal snapshot resume data: %w",
			err,
		)
	}

	err = json.Unmarshal(statusHistoryRaw, &view.StatusHistory)
	if err != nil {
		return nil, fmt.Errorf(
			"unmarshal status history: %w",
			err,
		)
	}

	return &view, nil
}
