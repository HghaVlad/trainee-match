package repository

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
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

func (a ApplicationRepo) Create(ctx context.Context, app application.Application) error {
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
