package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type CandidateRepo struct {
	db *pgxpool.Pool
}

func NewCandidateRepo(db *pgxpool.Pool) *CandidateRepo {
	return &CandidateRepo{db: db}
}

func (r *CandidateRepo) Create(ctx context.Context, candidate *domain.Candidate) (uuid.UUID, error) {
	query := `
		INSERT INTO candidates (id, user_id, full_name, phone, telegram, city, birthday) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query,
		candidate.ID,
		candidate.UserId,
		candidate.FullName,
		candidate.Phone,
		candidate.Telegram,
		candidate.City,
		candidate.Birthday,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "candidates_phone_key":
				return uuid.Nil, domain.ErrPhoneAlreadyExists
			case "candidates_telegram_key":
				return uuid.Nil, domain.ErrTelegramAlreadyExists
			case "candidates_user_id_key":
				return uuid.Nil, domain.ErrCandidateAlreadyExists
			}
		}
		return uuid.Nil, err
	}

	return id, nil
}

func (r *CandidateRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Candidate, error) {
	query := `
		SELECT id, user_id, full_name, phone, telegram, city, birthday 
		FROM candidates 
		WHERE id = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, id).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.FullName,
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
		SELECT id, user_id, full_name, phone, telegram, city, birthday 
		FROM candidates 
		WHERE user_id = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.FullName,
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
		SELECT id, user_id, full_name, phone, telegram, city, birthday 
		FROM candidates 
		WHERE telegram = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, telegram).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.FullName,
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
		SELECT id, user_id, full_name, phone, telegram, city, birthday 
		FROM candidates 
		WHERE phone = $1`

	var candidate domain.Candidate
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&candidate.ID,
		&candidate.UserId,
		&candidate.FullName,
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
		SET full_name = $1, phone = $2, telegram = $3, city = $4, birthday = $5 
		WHERE id = $6 
		RETURNING id, user_id, full_name, phone, telegram, city, birthday`

	var updatedCandidate domain.Candidate
	err := r.db.QueryRow(ctx, query,
		candidate.FullName,
		candidate.Phone,
		candidate.Telegram,
		candidate.City,
		candidate.Birthday,
		candidate.ID,
	).Scan(
		&updatedCandidate.ID,
		&updatedCandidate.UserId,
		&updatedCandidate.FullName,
		&updatedCandidate.Phone,
		&updatedCandidate.Telegram,
		&updatedCandidate.City,
		&updatedCandidate.Birthday,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Candidate{}, domain.ErrCandidateNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "candidates_phone_key":
				return domain.Candidate{}, domain.ErrPhoneAlreadyExists
			case "candidates_telegram_key":
				return domain.Candidate{}, domain.ErrTelegramAlreadyExists
			case "candidates_user_id_key":
				return domain.Candidate{}, domain.ErrCandidateAlreadyExists
			}
		}
		return domain.Candidate{}, err
	}

	return updatedCandidate, nil
}

func (r *CandidateRepo) GetCandidates(ctx context.Context, page int, size int) ([]domain.Candidate, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	offset := (page - 1) * size
	query := `
		SELECT id, user_id, full_name, phone, telegram, city, birthday
		FROM candidates
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candidates := make([]domain.Candidate, 0)
	for rows.Next() {
		var candidate domain.Candidate
		if err := rows.Scan(
			&candidate.ID,
			&candidate.UserId,
			&candidate.FullName,
			&candidate.Phone,
			&candidate.Telegram,
			&candidate.City,
			&candidate.Birthday,
		); err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}
