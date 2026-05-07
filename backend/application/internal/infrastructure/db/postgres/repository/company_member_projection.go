package repository

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
