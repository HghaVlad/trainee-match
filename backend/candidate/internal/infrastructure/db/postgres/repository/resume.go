package repository

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeRepo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewResumeRepo(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *ResumeRepo {
	return &ResumeRepo{db: db, getter: getter}
}

func (r *ResumeRepo) Create(ctx context.Context, resume *domain.Resume) (uuid.UUID, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	query := `INSERT INTO resumes (candidate_id, name, status, moderation_status, data) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id uuid.UUID
	err := conn.QueryRow(ctx, query, resume.CandidateId, resume.Name, resume.Status, resume.ModerationStatus, resume.Data).
		Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *ResumeRepo) GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error) {
	query := `SELECT id, candidate_id, name, status, moderation_status, data FROM resumes WHERE id = $1`

	var resume domain.Resume
	err := r.db.QueryRow(ctx, query, id).
		Scan(&resume.ID, &resume.CandidateId, &resume.Name, &resume.Status, &resume.ModerationStatus, &resume.Data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Resume{}, domain.ErrResumeNotFound
		}
		return domain.Resume{}, err
	}
	return resume, nil
}

func (r *ResumeRepo) GetByCandidateId(ctx context.Context, userId uuid.UUID, page, size int) ([]domain.Resume, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	offset := (page - 1) * size

	query := `SELECT id, candidate_id, name, status, moderation_status FROM resumes WHERE candidate_id = $1
				ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userId, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resumes = make([]domain.Resume, 0)
	for rows.Next() {
		var resume domain.Resume
		err = rows.Scan(&resume.ID, &resume.CandidateId, &resume.Name, &resume.Status, &resume.ModerationStatus)
		if err != nil {
			return nil, err
		}
		resumes = append(resumes, resume)
	}

	return resumes, nil
}

func (r *ResumeRepo) Update(ctx context.Context, resume *domain.Resume) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	query := `UPDATE resumes SET name = $1, status = $2, data = $3 WHERE id = $4`

	cmdTag, err := conn.Exec(ctx, query, resume.Name, resume.Status, resume.Data, resume.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrResumeNotFound
	}
	return nil
}

func (r *ResumeRepo) Remove(ctx context.Context, id, candidateID uuid.UUID) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	query := `DELETE FROM resumes WHERE id = $1 AND candidate_id = $2`

	result, err := conn.Exec(ctx, query, id, candidateID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrResumeNotFound
	}

	return nil
}

func (r *ResumeRepo) SetModerationStatus(ctx context.Context, id uuid.UUID, status domain.ModerationStatus) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	query := `UPDATE resumes SET moderation_status = $1 WHERE id = $2`

	result, err := conn.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrResumeNotFound
	}

	return nil
}
