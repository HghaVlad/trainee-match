package repository

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type VacancyProjection struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewVacancyProjection(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *VacancyProjection {
	return &VacancyProjection{
		db:     db,
		getter: getter,
	}
}

// TODO: check in apply uc

func (p *VacancyProjection) GetByID(ctx context.Context, vacID uuid.UUID) (*projection.Vacancy, error) {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		SELECT id, company_id, company_name, title, status, created_at, 
		       updated_at, moderation_status, company_moderation_status 
		FROM vacancy_projection
		WHERE id = $1
	`

	var vacancy projection.Vacancy

	err := q.QueryRow(ctx, query, vacID).
		Scan(&vacancy.ID, &vacancy.CompanyID, &vacancy.CompanyName, &vacancy.Title,
			&vacancy.Status, &vacancy.CreatedAt, &vacancy.UpdatedAt,
			&vacancy.ModStatus, &vacancy.CompModStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, projection.ErrVacancyNotFound
		}

		return nil, fmt.Errorf("get vacancy projection: %w", err)
	}

	return &vacancy, nil
}

// CheckHrAccess returns companyID if success,
// projection.ErrVacancyNotFound otherwise
func (p *VacancyProjection) CheckHrAccess(
	ctx context.Context,
	userID, vacancyID uuid.UUID,
) (uuid.UUID, error) {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `SELECT v.company_id
		FROM vacancy_projection v
		JOIN company_members cm ON cm.company_id = v.company_id
		WHERE v.id = $1 AND cm.user_id = $2`

	var companyID uuid.UUID

	err := q.QueryRow(ctx, query, vacancyID, userID).Scan(&companyID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, projection.ErrVacancyNotFound
		}

		return uuid.Nil, fmt.Errorf("check hr access company id: %w", err)
	}

	return companyID, nil
}

func (p *VacancyProjection) Save(ctx context.Context, vacancy *projection.Vacancy) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		INSERT INTO vacancy_projection
			(id, company_id, company_name, title, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			company_id = EXCLUDED.company_id,
			company_name = EXCLUDED.company_name,
			title = EXCLUDED.title,
			status = EXCLUDED.status,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := q.Exec(ctx, query,
		vacancy.ID,
		vacancy.CompanyID,
		vacancy.CompanyName,
		vacancy.Title,
		vacancy.Status,
		vacancy.CreatedAt,
		vacancy.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save vacancy projection: %w", err)
	}
	return nil
}

func (p *VacancyProjection) UpdateTitle(ctx context.Context, vacancyID uuid.UUID, title string) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		UPDATE vacancy_projection
		SET title = $2
		WHERE id = $1
	`

	_, err := q.Exec(ctx, query, vacancyID, title)
	if err != nil {
		return fmt.Errorf("update vacancy projection title: %w", err)
	}
	return nil
}

func (p *VacancyProjection) Archive(ctx context.Context, vacancyID uuid.UUID) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		UPDATE vacancy_projection
		SET status = $2
		WHERE id = $1
	`

	_, err := q.Exec(ctx, query, vacancyID, projection.VacancyStatusArchived)
	if err != nil {
		return fmt.Errorf("archive vacancy projection: %w", err)
	}
	return nil
}

func (p *VacancyProjection) UpdateCompanyName(ctx context.Context, companyID uuid.UUID, companyName string) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		UPDATE vacancy_projection
		SET company_name = $2
		WHERE company_id = $1
	`

	_, err := q.Exec(ctx, query, companyID, companyName)
	if err != nil {
		return fmt.Errorf("update vacancy projection company name: %w", err)
	}
	return nil
}

func (p *VacancyProjection) ArchiveByCompany(ctx context.Context, compID uuid.UUID) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		UPDATE vacancy_projection
		SET status = $2
		WHERE company_id = $1`

	_, err := q.Exec(ctx, query, compID, projection.VacancyStatusArchived)

	if err != nil {
		return fmt.Errorf("archive vacancies by company: %w", err)
	}

	return nil
}

func (p *VacancyProjection) UpdateModStatus(
	ctx context.Context,
	vacID uuid.UUID,
	status projection.ModerationStatus,
) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		UPDATE vacancy_projection
		SET moderation_status = $2
		WHERE id = $1`

	_, err := q.Exec(ctx, query, vacID, status)

	if err != nil {
		return fmt.Errorf("vacancy upd mod status: %w", err)
	}

	return nil
}

func (p *VacancyProjection) UpdateCompanyModStatus(
	ctx context.Context,
	compID uuid.UUID,
	status projection.ModerationStatus,
) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `
		UPDATE vacancy_projection
		SET company_moderation_status = $2
		WHERE company_id = $1`

	_, err := q.Exec(ctx, query, compID, status)

	if err != nil {
		return fmt.Errorf("vacancy upd company mod status: %w", err)
	}

	return nil
}

func (p *VacancyProjection) DeleteByCompanyID(ctx context.Context, companyID uuid.UUID) error {
	q := p.getter.DefaultTrOrDB(ctx, p.db)

	const query = `DELETE FROM vacancy_projection WHERE company_id = $1`

	_, err := q.Exec(ctx, query, companyID)
	if err != nil {
		return fmt.Errorf("delete vacancy projection by company: %w", err)
	}
	return nil
}
