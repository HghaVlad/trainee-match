package repository

import (
	"context"
	"errors"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CandidateRepo struct {
	db *pgxpool.Pool
}

func NewCandidateRepo(db *pgxpool.Pool) *CandidateRepo {
	return &CandidateRepo{db: db}
}

func (r *CandidateRepo) Create(ctx context.Context, candidate *domain.Candidate) (uuid.UUID, error) {
	query := `
		INSERT INTO candidates (user_id, phone, telegram, city, birthday) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query,
		candidate.UserId,
		candidate.Phone,
		candidate.Telegram,
		candidate.City,
		candidate.Birthday,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *CandidateRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Candidate, error) {
	query := `
		SELECT id, user_id, phone, telegram, city, birthday 
		FROM candidates 
		WHERE id = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, id).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.Phone,
		&candidate.Telegram,
		&candidate.City,
		&candidate.Birthday,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Candidate{}, domain.ErrCandidateNotFound
		}
		return domain.Candidate{}, err
	}

	return candidate, nil
}

func (r *CandidateRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (domain.Candidate, error) {
	query := `
		SELECT id, user_id, phone, telegram, city, birthday 
		FROM candidates 
		WHERE user_id = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.Phone,
		&candidate.Telegram,
		&candidate.City,
		&candidate.Birthday,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Candidate{}, domain.ErrCandidateNotFound
		}
		return domain.Candidate{}, err
	}

	return candidate, nil
}

func (r *CandidateRepo) GetByTelegram(ctx context.Context, telegram string) (domain.Candidate, error) {
	query := `
		SELECT id, user_id, phone, telegram, city, birthday 
		FROM candidates 
		WHERE telegram = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, telegram).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.Phone,
		&candidate.Telegram,
		&candidate.City,
		&candidate.Birthday,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Candidate{}, domain.ErrCandidateNotFound
		}
		return domain.Candidate{}, err
	}

	return candidate, nil
}

func (r *CandidateRepo) GetByPhone(ctx context.Context, phone string) (domain.Candidate, error) {
	query := `
		SELECT id, user_id, phone, telegram, city, birthday 
		FROM candidates 
		WHERE phone = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.Phone,
		&candidate.Telegram,
		&candidate.City,
		&candidate.Birthday,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Candidate{}, domain.ErrCandidateNotFound
		}
		return domain.Candidate{}, err
	}

	return candidate, nil
}

func (r *CandidateRepo) Update(ctx context.Context, candidate domain.Candidate) (domain.Candidate, error) {
	query := `
		UPDATE candidates 
		SET phone = $1, telegram = $2, city = $3, birthday = $4 
		WHERE id = $5 
		RETURNING id, user_id, phone, telegram, city, birthday`

	var updatedCandidate domain.Candidate
	err := r.db.QueryRow(ctx, query,
		candidate.Phone,
		candidate.Telegram,
		candidate.City,
		candidate.Birthday,
		candidate.ID,
	).Scan(
		&updatedCandidate.ID,
		&updatedCandidate.UserId,
		&updatedCandidate.Phone,
		&updatedCandidate.Telegram,
		&updatedCandidate.City,
		&updatedCandidate.Birthday,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Candidate{}, domain.ErrCandidateNotFound
		}
		return domain.Candidate{}, err
	}

	return updatedCandidate, nil
}
