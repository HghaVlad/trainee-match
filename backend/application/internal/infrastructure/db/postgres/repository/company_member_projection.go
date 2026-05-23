package repository

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type CompanyMemberProjection struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewCompanyMemberProjection(db *pgxpool.Pool, getter *trmpgx.CtxGetter) *CompanyMemberProjection {
	return &CompanyMemberProjection{
		db:     db,
		getter: getter,
	}
}

func (r *CompanyMemberProjection) IsMember(
	ctx context.Context,
	userID uuid.UUID,
	companyID uuid.UUID,
) (bool, error) {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `SELECT EXISTS (
			SELECT 1
			FROM company_members
			WHERE user_id = $1 AND company_id = $2
		)`

	var exists bool

	err := q.QueryRow(ctx, query, userID, companyID).
		Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check company membership: %w", err)
	}

	return exists, nil
}

func (r *CompanyMemberProjection) Save(ctx context.Context, member *projection.CompanyMember) error {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `
		INSERT INTO company_members (user_id, company_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, company_id) DO UPDATE SET
			role = EXCLUDED.role
	`

	_, err := q.Exec(ctx, query, member.UserID, member.CompanyID, member.Role)
	if err != nil {
		return fmt.Errorf("save company member projection: %w", err)
	}
	return nil
}

func (r *CompanyMemberProjection) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `DELETE FROM company_members WHERE user_id = $1 AND company_id = $2`

	_, err := q.Exec(ctx, query, userID, companyID)
	if err != nil {
		return fmt.Errorf("delete company member projection: %w", err)
	}
	return nil
}

func (r *CompanyMemberProjection) DeleteByCompanyID(ctx context.Context, companyID uuid.UUID) error {
	q := r.getter.DefaultTrOrDB(ctx, r.db)

	const query = `DELETE FROM company_members WHERE company_id = $1`

	_, err := q.Exec(ctx, query, companyID)
	if err != nil {
		return fmt.Errorf("delete company members by company: %w", err)
	}
	return nil
}
