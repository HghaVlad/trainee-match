package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/summary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
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

		return fmt.Errorf("create application: %w", err)
	}

	return nil
}

func (a *ApplicationRepo) GetCandidateDetailedView(
	ctx context.Context,
	appID, candID uuid.UUID,
) (*views.CandidateDetailedView, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `
		SELECT a.id, a.status, a.cover_letter, a.created_at, a.updated_at,
        	v.id, v.company_id, v.company_name, v.title,
        	s.resume_data, s.email, s.full_name, s.telegram, s.created_at,
        	COALESCE(h.status_history, '[]') AS status_history

		FROM applications a
		LEFT JOIN vacancy_projection v ON v.id = a.vacancy_id
		JOIN application_snapshots s ON a.snapshot_id = s.id
		
		LEFT JOIN LATERAL (
		    SELECT json_agg(
		    	json_build_object(
					'status', h.status,
					'created_at', h.created_at, 
					'changed_by_role', h.changed_by_role
				)
				ORDER BY h.created_at 
			) AS status_history
		    FROM application_status_history h
    		WHERE h.application_id = a.id
		) h ON TRUE
		
 		WHERE a.id = $1 AND a.candidate_id = $2`

	row := q.QueryRow(ctx, query, appID, candID)

	view, err := scanCandidateDetailedView(row)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrNotFound
		}

		return nil, fmt.Errorf("get candidate detailed app view: %w", err)
	}

	return view, nil
}

func (a *ApplicationRepo) GetForUpdateByCandidate(
	ctx context.Context,
	appID, candID uuid.UUID,
) (*application.Application, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `
		SELECT id, resume_id, candidate_id, vacancy_id, company_id,
       		snapshot_id, status, cover_letter, created_at, updated_at
		FROM applications
		WHERE id = $1 AND candidate_id = $2
		FOR UPDATE`

	var app application.Application

	row := q.QueryRow(ctx, query, appID, candID)
	err := scanApp(row, &app)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrNotFound
		}

		return nil, fmt.Errorf("get application: %w", err)
	}

	return &app, nil
}

func (a *ApplicationRepo) GetHrDetailedView(
	ctx context.Context,
	appID, hrID uuid.UUID,
) (*views.HrDetailedView, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `
		SELECT a.id, a.status, a.cover_letter, a.created_at, a.updated_at,
        		v.id, v.title,
        		s.resume_data, s.email, s.full_name, s.telegram, s.created_at,
        		COALESCE(h.status_history, '[]') AS status_history
		
		FROM applications a
		LEFT JOIN vacancy_projection v ON v.id = a.vacancy_id
		JOIN application_snapshots s ON a.snapshot_id = s.id
		JOIN company_members cm ON cm.company_id = a.company_id
		
		LEFT JOIN LATERAL (
		    SELECT json_agg(
		    	json_build_object(
					'status', h.status,
					'created_at', h.created_at,
					'comment', h.comment,  
					'changed_by_role', h.changed_by_role,
					'changed_by_user_id', h.changed_by_user_id
				)
				ORDER BY h.created_at 
			) AS status_history
		    FROM application_status_history h
    		WHERE h.application_id = a.id
		) h ON TRUE
		
		WHERE a.id = $1 AND cm.user_id = $2`

	row := q.QueryRow(ctx, query, appID, hrID)

	view, err := scanHrDetailedView(row)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrNotFound
		}

		return nil, fmt.Errorf("get hr detailed app view: %w", err)
	}

	return view, nil
}

func (a *ApplicationRepo) GetForUpdateByHr(
	ctx context.Context,
	appID, hrID uuid.UUID,
) (*application.Application, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `
		SELECT a.id, a.resume_id, a.candidate_id, a.vacancy_id, a.company_id,
       		a.snapshot_id, a.status, a.cover_letter, a.created_at, a.updated_at
		FROM applications a
		JOIN company_members cm ON cm.company_id = a.company_id
		WHERE a.id = $1 AND cm.user_id = $2
		FOR UPDATE`

	var app application.Application

	row := q.QueryRow(ctx, query, appID, hrID)
	err := scanApp(row, &app)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, application.ErrNotFound
		}

		return nil, fmt.Errorf("get application: %w", err)
	}

	return &app, nil
}

func (a *ApplicationRepo) UpdateStatus(
	ctx context.Context,
	appID uuid.UUID,
	status application.Status,
	updAt time.Time,
) error {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	const query = `UPDATE applications
		SET status = $1, updated_at = $2
		WHERE id = $3`

	cmdTag, err := q.Exec(ctx, query, status, updAt, appID)
	if err != nil {
		return fmt.Errorf("update application status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return application.ErrNotFound
	}

	return nil
}

func (a *ApplicationRepo) ListCandidateAppSummaries(
	ctx context.Context,
	candidateID uuid.UUID,
	statuses []application.Status,
	companyID *uuid.UUID,
	order cursors.SummaryOrder,
	cursor any,
	limit int,
) ([]views.CandidateAppSummary, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	orderByColumn := "a.created_at"
	if order == cursors.OrderUpdatedAtDesc {
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

	if cur, ok := cursor.(*cursors.SummaryCursor); ok && cur != nil {
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

	items := make([]views.CandidateAppSummary, 0, limit)
	for rows.Next() {
		var item views.CandidateAppSummary

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

func (a *ApplicationRepo) ListHrAppSummaries(
	ctx context.Context,
	hrUserID uuid.UUID,
	statuses []application.Status,
	companyID, vacancyID *uuid.UUID,
	createdFrom, createdTo *time.Time,
	order cursors.HrSummaryOrder,
	cursor any,
	limit int,
) ([]views.HrAppSummary, error) {
	q := a.getter.DefaultTrOrDB(ctx, a.db)

	b := buildHrSummaryListQuery(
		hrUserID,
		statuses,
		companyID,
		vacancyID,
		createdFrom,
		createdTo,
		order,
		cursor,
		limit,
	)

	query := fmt.Sprintf(`
		SELECT a.id, a.status, a.vacancy_id, COALESCE(v.title, ''),
			s.email, s.full_name, s.telegram, s.created_at,
			a.created_at, a.updated_at
		FROM applications a
		JOIN company_members cm ON cm.company_id = a.company_id
		LEFT JOIN vacancy_projection v ON v.id = a.vacancy_id
		JOIN application_snapshots s ON s.id = a.snapshot_id
		WHERE %s
		ORDER BY %s
		LIMIT $%d
	`,
		strings.Join(b.conditions, " AND "),
		b.orderBy,
		b.limitPos,
	)

	rows, err := q.Query(ctx, query, b.args...)
	if err != nil {
		return nil, fmt.Errorf("list hr summaries: %w", err)
	}
	defer rows.Close()

	items := make([]views.HrAppSummary, 0, limit)
	for rows.Next() {
		var item views.HrAppSummary

		err := rows.Scan(&item.AppID, &item.Status, &item.VacancyID, &item.VacancyTitle,
			&item.AppSnap.Email, &item.AppSnap.FullName, &item.AppSnap.Telegram,
			&item.AppSnap.CreatedAt, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan hr summary: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate hr summaries: %w", err)
	}

	return items, nil
}

func (a *ApplicationRepo) GetCompanyAnalyticsSummary(
	ctx context.Context,
	compID uuid.UUID,
) (*summary.Summary, error) {
	const query = `
		SELECT COUNT(*) FILTER (WHERE status = 'submitted') AS submitted,
		    COUNT(*) FILTER (WHERE status = 'seen') AS seen,
			COUNT(*) FILTER (WHERE status = 'interview') AS interview,
			COUNT(*) FILTER (WHERE status = 'offer') AS offer,
			COUNT(*) FILTER (WHERE status = 'rejected') AS rejected,
			COUNT(*) FILTER (WHERE status = 'withdrawn') AS withdrawn
		FROM applications
		WHERE company_id = $1`

	var sum summary.Summary
	row := a.db.QueryRow(ctx, query, compID)
	err := scanSummary(row, &sum)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, projection.ErrCompanyNotFound
		}

		return nil, fmt.Errorf("get company analytics summary: %w", err)
	}

	return &sum, nil
}

func (a *ApplicationRepo) GetVacancyAnalyticsSummary(
	ctx context.Context,
	vacID uuid.UUID,
) (*summary.Summary, error) {
	const query = `
		SELECT COUNT(*) FILTER (WHERE status = 'submitted') AS submitted,
		    COUNT(*) FILTER (WHERE status = 'seen') AS seen,
			COUNT(*) FILTER (WHERE status = 'interview') AS interview,
			COUNT(*) FILTER (WHERE status = 'offer') AS offer,
			COUNT(*) FILTER (WHERE status = 'rejected') AS rejected,
			COUNT(*) FILTER (WHERE status = 'withdrawn') AS withdrawn
		FROM applications
		WHERE vacancy_id = $1`

	var sum summary.Summary
	row := a.db.QueryRow(ctx, query, vacID)
	err := scanSummary(row, &sum)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, projection.ErrVacancyNotFound
		}

		return nil, fmt.Errorf("get vacancy analytics summary: %w", err)
	}

	return &sum, nil
}

func scanCandidateDetailedView(row pgx.Row) (*views.CandidateDetailedView, error) {
	var (
		view             views.CandidateDetailedView
		resumeDataRaw    []byte
		statusHistoryRaw []byte
	)

	err := row.Scan(&view.AppID, &view.Status, &view.CoverLetter, &view.CreatedAt, &view.UpdatedAt,
		&view.VacancyID, &view.CompanyID, &view.CompanyName, &view.VacancyTitle,
		&resumeDataRaw, &view.Snapshot.Email, &view.Snapshot.FullName, &view.Snapshot.Telegram,
		&view.Snapshot.CreatedAt, &statusHistoryRaw)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resumeDataRaw, &view.Snapshot.ResumeData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal snapshot resume data: %w", err)
	}

	err = json.Unmarshal(statusHistoryRaw, &view.StatusHistory)
	if err != nil {
		return nil, fmt.Errorf("unmarshal status history: %w", err)
	}

	if view.StatusHistory == nil {
		view.StatusHistory = []views.StatusChangeCandidateView{}
	}

	return &view, nil
}

func scanHrDetailedView(row pgx.Row) (*views.HrDetailedView, error) {
	var (
		view             views.HrDetailedView
		resumeDataRaw    []byte
		statusHistoryRaw []byte
	)

	err := row.Scan(&view.AppID, &view.Status, &view.CoverLetter, &view.CreatedAt, &view.UpdatedAt,
		&view.VacancyID, &view.VacancyTitle,
		&resumeDataRaw, &view.Snapshot.Email, &view.Snapshot.FullName,
		&view.Snapshot.Telegram, &view.Snapshot.CreatedAt,
		&statusHistoryRaw,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resumeDataRaw, &view.Snapshot.ResumeData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal snapshot resume data: %w", err)
	}

	err = json.Unmarshal(statusHistoryRaw, &view.StatusHistory)
	if err != nil {
		return nil, fmt.Errorf("unmarshal status history: %w", err)
	}

	if view.StatusHistory == nil {
		view.StatusHistory = []views.StatusChangeHrFullView{}
	}

	return &view, nil
}

func scanApp(row pgx.Row, app *application.Application) error {
	return row.Scan(&app.ID, &app.ResumeID, &app.CandidateID, &app.VacancyID, &app.CompanyID,
		&app.SnapshotID, &app.Status, &app.CoverLetter, &app.CreatedAt, &app.UpdatedAt)
}

func scanSummary(row pgx.Row, sum *summary.Summary) error {
	return row.Scan(&sum.SubmittedCount, &sum.SeenCount, &sum.InterviewCount,
		&sum.OfferCount, &sum.RejectedCount, &sum.WithdrawnCount)
}

type hrSummaryListQuery struct {
	args       []any
	conditions []string
	orderBy    string
	limitPos   int
}

func buildHrSummaryListQuery(
	hrUserID uuid.UUID,
	statuses []application.Status,
	companyID *uuid.UUID,
	vacancyID *uuid.UUID,
	createdFrom *time.Time,
	createdTo *time.Time,
	order cursors.HrSummaryOrder,
	cursor any,
	limit int,
) hrSummaryListQuery {
	q := hrSummaryListQuery{
		args: []any{hrUserID},
		conditions: []string{
			"cm.user_id = $1",
			"cm.company_id = a.company_id",
		},
		orderBy: "a.created_at DESC, a.id DESC",
	}

	if len(statuses) > 0 {
		strStatuses := make([]string, 0, len(statuses))
		for _, s := range statuses {
			strStatuses = append(strStatuses, string(s))
		}

		q.args = append(q.args, strStatuses)
		q.conditions = append(
			q.conditions,
			fmt.Sprintf("a.status = ANY($%d::application_status_enum[])", len(q.args)),
		)
	}

	if companyID != nil {
		q.args = append(q.args, *companyID)
		q.conditions = append(
			q.conditions,
			fmt.Sprintf("a.company_id = $%d", len(q.args)),
		)
	}

	if vacancyID != nil {
		q.args = append(q.args, *vacancyID)
		q.conditions = append(
			q.conditions,
			fmt.Sprintf("a.vacancy_id = $%d", len(q.args)),
		)
	}

	if createdFrom != nil {
		q.args = append(q.args, *createdFrom)
		q.conditions = append(
			q.conditions,
			fmt.Sprintf("a.created_at >= $%d", len(q.args)),
		)
	}

	if createdTo != nil {
		q.args = append(q.args, *createdTo)
		q.conditions = append(
			q.conditions,
			fmt.Sprintf("a.created_at <= $%d", len(q.args)),
		)
	}

	switch order {
	case cursors.HrSummaryOrderUpdatedAtDesc:
		q.orderBy = "a.updated_at DESC, a.id DESC"

	case cursors.HrSummaryOrderCandidateFullName:
		q.orderBy = "s.full_name ASC, a.id ASC"
	}

	applyHrSummaryCursor(&q, order, cursor)

	q.args = append(q.args, limit)
	q.limitPos = len(q.args)

	return q
}

func applyHrSummaryCursor(
	q *hrSummaryListQuery,
	order cursors.HrSummaryOrder,
	cursor any,
) {
	cur, ok := cursor.(*cursors.HrSummaryCursor)
	if !ok || cur == nil {
		return
	}

	switch order {
	case cursors.HrSummaryOrderUpdatedAtDesc:
		applyUpdatedAtCursor(q, cur)

	case cursors.HrSummaryOrderCandidateFullName:
		applyCandidateNameCursor(q, cur)

	default:
		applyCreatedAtCursor(q, cur)
	}
}

func applyUpdatedAtCursor(
	q *hrSummaryListQuery,
	cur *cursors.HrSummaryCursor,
) {
	if cur.SortAt == nil {
		return
	}

	q.args = append(q.args, *cur.SortAt, cur.AppID)

	sortPos := len(q.args) - 1
	idPos := len(q.args)

	q.conditions = append(q.conditions,
		fmt.Sprintf("(a.updated_at < $%d OR (a.updated_at = $%d AND a.id < $%d))",
			sortPos, sortPos, idPos,
		),
	)
}

func applyCreatedAtCursor(
	q *hrSummaryListQuery,
	cur *cursors.HrSummaryCursor,
) {
	if cur.SortAt == nil {
		return
	}

	q.args = append(q.args, *cur.SortAt, cur.AppID)

	sortPos := len(q.args) - 1
	idPos := len(q.args)

	q.conditions = append(q.conditions,
		fmt.Sprintf("(a.created_at < $%d OR (a.created_at = $%d AND a.id < $%d))",
			sortPos, sortPos, idPos,
		),
	)
}

func applyCandidateNameCursor(
	q *hrSummaryListQuery,
	cur *cursors.HrSummaryCursor,
) {
	if cur.FullName == nil {
		return
	}

	q.args = append(q.args, *cur.FullName, cur.AppID)

	namePos := len(q.args) - 1
	idPos := len(q.args)

	q.conditions = append(q.conditions,
		fmt.Sprintf("(s.full_name > $%d OR (s.full_name = $%d AND a.id > $%d))",
			namePos, namePos, idPos,
		),
	)
}
