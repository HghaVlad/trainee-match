package repository

import (
	"context"
	"errors"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type SkillRepo struct {
	db *pgxpool.Pool
}

func NewSkillRepo(db *pgxpool.Pool) *SkillRepo {
	return &SkillRepo{db: db}
}

func (r *SkillRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Skill, error) {
	query := `SELECT id, name FROM skills WHERE id = $1`
	var skill domain.Skill
	err := r.db.QueryRow(ctx, query, id).Scan(&skill.ID, &skill.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Skill{}, domain.ErrResumeNotFound
		}
		return domain.Skill{}, err
	}
	return skill, err
}

func (r *SkillRepo) List(ctx context.Context) ([]domain.Skill, error) {
	query := `SELECT id, name FROM skills ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.Skill
	for rows.Next() {
		var skill domain.Skill
		if err := rows.Scan(&skill.ID, &skill.Name); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (r *SkillRepo) AreSkillsExist(ctx context.Context, skillIDs []uuid.UUID) (bool, error) {
	if len(skillIDs) == 0 {
		return true, nil
	}

	query := `SELECT COUNT(*) FROM skills WHERE id = ANY($1)`
	var count int
	err := r.db.QueryRow(ctx, query, pq.Array(skillIDs)).Scan(&count)
	if err != nil {
		return false, err
	}

	// If count equals the length of skillIDs, all skills exist
	return count == len(skillIDs), nil
}
