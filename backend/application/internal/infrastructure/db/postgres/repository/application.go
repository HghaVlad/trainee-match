package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	listcandidatesummary "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
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

func (a *ApplicationRepo) ListCandidateSummaries(
	ctx context.Context,
	candidateID uuid.UUID,
	statuses []application.Status,
	companyID *uuid.UUID,
	order listcandidatesummary.Order,
	cursor any,
	limit int,
) ([]views.CandidateSummary, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	orderByColumn := "a.created_at"
	if order == listcandidatesummary.OrderUpdatedAtDesc {
		orderByColumn = "a.updated_at"
	}

	args := []any{candidateID}
	conditions := []string{"a.candidate_id = $1"}

	if len(statuses) > 0 {
		strStatuses := make([]string, 0, len(statuses))
		for _, s := range statuses {
			strStatuses = append(strStatuses, string(s))
		}

		args = append(args, strStatuses)
		conditions = append(conditions, fmt.Sprintf("a.status = ANY($%d::application_status_enum[])", len(args)))
	}

	if companyID != nil {
		args = append(args, *companyID)
		conditions = append(conditions, fmt.Sprintf("a.company_id = $%d", len(args)))
	}

	if cur, ok := cursor.(*listcandidatesummary.SummaryCursor); ok && cur != nil {
		args = append(args, cur.SortAt, cur.AppID)
		sortArgPos := len(args) - 1
		idArgPos := len(args)
		conditions = append(conditions,
			fmt.Sprintf("(%s < $%d OR (%s = $%d AND a.id < $%d))",
				orderByColumn, sortArgPos, orderByColumn, sortArgPos, idArgPos),
		)
	}

	args = append(args, limit)
	limitPos := len(args)

	query := fmt.Sprintf(`
		SELECT a.id, a.status, a.vacancy_id, COALESCE(v.title, ''), a.company_id, COALESCE(v.company_name, ''), a.created_at, a.updated_at
		FROM applications a
		LEFT JOIN vacancy_projection v ON v.id = a.vacancy_id
		WHERE %s
		ORDER BY %s DESC, a.id DESC
		LIMIT $%d
	`, strings.Join(conditions, " AND "), orderByColumn, limitPos)

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list candidate summaries: %w", err)
	}
	defer rows.Close()

	items := make([]views.CandidateSummary, 0, limit)
	for rows.Next() {
		var item views.CandidateSummary

		err := rows.Scan(&item.AppID, &item.Status, &item.VacancyID, &item.VacancyTitle,
			&item.CompanyID, &item.CompanyName, &item.CreatedAt, &item.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan candidate summary: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate candidate summaries: %w", err)
	}

	return items, nil
}
