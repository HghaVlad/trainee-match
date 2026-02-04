package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *SkillRepo) CheckExistsBatch(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]bool{}, nil
	}
	
	// Create placeholders for the query
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	
	query := fmt.Sprintf("SELECT id FROM skills WHERE id IN (%s)", strings.Join(placeholders, ", "))
	
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	existingIds := make(map[uuid.UUID]bool)
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		existingIds[id] = true
	}
	
	return existingIds, nil
}
