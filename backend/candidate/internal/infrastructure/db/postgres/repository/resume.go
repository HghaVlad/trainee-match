package repository

import (
	"context"
	"errors"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResumeRepo struct {
	db *pgxpool.Pool
}

func NewResumeRepo(db *pgxpool.Pool) *ResumeRepo {
	return &ResumeRepo{db: db}
}

func (r *ResumeRepo) Create(ctx context.Context, resume *domain.Resume) (uuid.UUID, error) {
	query := `INSERT INTO resume (candidate_id, name, status, data) VALUES ($1, $2, $3, $4) RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, resume.CandidateId, resume.Name, resume.Status, resume.Data).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *ResumeRepo) GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error) {
	query := `SELECT id, candidate_id, name, status, data FROM resume WHERE id = $1`

	var resume domain.Resume
	err := r.db.QueryRow(ctx, query, id).Scan(&resume.ID, &resume.CandidateId, &resume.Name, &resume.Status, &resume.Data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Resume{}, domain.ErrResumeNotFound
		}
		return domain.Resume{}, err
	}
	return resume, nil
}

func (r *ResumeRepo) GetByCandidateId(ctx context.Context, userId uuid.UUID) ([]domain.Resume, error) {
	query := `SELECT id, candidate_id, name, status FROM resume WHERE candidate_id = $1`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resumes = make([]domain.Resume, 0)
	for rows.Next() {
		var resume domain.Resume
		err = rows.Scan(&resume.ID, &resume.CandidateId, &resume.Name, &resume.Status)
		if err != nil {
			return nil, err
		}
		resumes = append(resumes, resume)
	}

	return resumes, nil
}

func (r *ResumeRepo) Update(ctx context.Context, resume *domain.Resume) error {
	query := `UPDATE resume SET name = $1, status = $2, data = $3 WHERE id = $4`

	_, err := r.db.Exec(ctx, query, resume.Name, resume.Status, resume.Data, resume.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrResumeNotFound
		}
		return err
	}
	return nil
}

func (r *ResumeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM resume WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrResumeNotFound
	}

	return nil
}
