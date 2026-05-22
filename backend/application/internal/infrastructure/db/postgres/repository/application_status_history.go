package repository

import (
	"context"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/dynamics"
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
			WHERE application_id = $1 AND a.candidate_id = $2
			ORDER BY sh.created_at, sh.id`

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
		WHERE application_id = $1 AND cm.user_id = $2
		ORDER BY sh.created_at, sh.id`

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

func (r *AppStatusHistoryRepo) AddChangesByVacancy(
	ctx context.Context,
	vacID uuid.UUID,
	newStatus application.Status,
	statusesToUpdate []application.Status,
	role application.Actor,
	comment *string,
	when time.Time,
) error {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
	INSERT INTO application_status_history (
		application_id, status, changed_by_role, comment, created_at
	)
	SELECT id, $2, $3, $4, $5
	FROM applications
	WHERE vacancy_id = $1 AND status = ANY($6)`

	_, err := q.Exec(ctx, query, vacID, newStatus, role, comment, when, appStatusesToStrings(statusesToUpdate))

	if err != nil {
		return fmt.Errorf("add app changes be vacancy: %w", err)
	}

	return nil
}

func (r *AppStatusHistoryRepo) GetDynamicsBucketsByCompany(
	ctx context.Context,
	compID uuid.UUID,
	from, to time.Time,
	interval dynamics.Interval,
) ([]dynamics.Bucket, error) {
	const query = `
		WITH buckets AS (
			SELECT generate_series($2::timestamptz, $3::timestamptz, $4::interval) AS bucket_start
		)	
		SELECT
			b.bucket_start, 
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'submitted'), 0) AS submitted,
		    COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'seen'), 0) AS seen,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'interview'), 0) AS interview,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'offer'), 0) AS offer,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'rejected'), 0) AS rejected,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'withdrawn'), 0) AS withdrawn
		FROM buckets b
		LEFT JOIN (
    		application_status_history ash
    		JOIN applications a ON a.id = ash.application_id AND a.company_id = $1
		)
		ON date_trunc('%s', ash.created_at)::timestamptz = b.bucket_start
		GROUP BY b.bucket_start
		ORDER BY b.bucket_start`

	prepared := fmt.Sprintf(query, interval)
	intervalStr := fmt.Sprintf("1 %s", interval)

	rows, err := r.db.Query(ctx, prepared, compID, from, to, intervalStr)
	if err != nil {
		return nil, fmt.Errorf("get dynamics buckets by company: %w", err)
	}

	defer rows.Close()
	buckets, err := scanBuckets(rows)
	if err != nil {
		return nil, fmt.Errorf("get dynamics buckets by company: %w", err)
	}

	return buckets, nil
}

func (r *AppStatusHistoryRepo) GetDynamicsBucketsByVacancy(
	ctx context.Context,
	vacID uuid.UUID,
	from, to time.Time,
	interval dynamics.Interval,
) ([]dynamics.Bucket, error) {
	const query = `
		WITH buckets AS (
			SELECT generate_series($2::timestamptz, $3::timestamptz, $4::interval) AS bucket_start
		)	
		SELECT
			b.bucket_start, 
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'submitted'), 0) AS submitted,
		    COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'seen'), 0) AS seen,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'interview'), 0) AS interview,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'offer'), 0) AS offer,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'rejected'), 0) AS rejected,
			COALESCE(COUNT(ash.id) FILTER (WHERE ash.status = 'withdrawn'), 0) AS withdrawn
		FROM buckets b
		LEFT JOIN (
    		application_status_history ash
    		JOIN applications a ON a.id = ash.application_id AND a.vacancy_id = $1
		)
		ON date_trunc('%s', ash.created_at)::timestamptz = b.bucket_start
		GROUP BY b.bucket_start
		ORDER BY b.bucket_start`

	prepared := fmt.Sprintf(query, interval)
	intervalStr := fmt.Sprintf("1 %s", interval)

	rows, err := r.db.Query(ctx, prepared, vacID, from, to, intervalStr)
	if err != nil {
		return nil, fmt.Errorf("get dynamics buckets by vacancy: %w", err)
	}

	defer rows.Close()
	buckets, err := scanBuckets(rows)
	if err != nil {
		return nil, fmt.Errorf("get dynamics buckets by vacancy: %w", err)
	}

	return buckets, nil
}

func scanBuckets(rows pgx.Rows) ([]dynamics.Bucket, error) {
	var buckets []dynamics.Bucket

	for rows.Next() {
		var bucket dynamics.Bucket
		err := rows.Scan(&bucket.Start, &bucket.SubmittedCount, &bucket.SeenCount,
			&bucket.InterviewCount, &bucket.OfferCount, &bucket.RejectedCount,
			&bucket.WithdrawnCount)
		if err != nil {
			return nil, fmt.Errorf("scan dynamics bucket: %w", err)
		}

		buckets = append(buckets, bucket)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan dynamics bucket: %w", err)
	}

	return buckets, nil
}

func appStatusesToStrings(status []application.Status) []string {
	strs := make([]string, len(status))
	for i, s := range status {
		strs[i] = string(s)
	}
	return strs
}
