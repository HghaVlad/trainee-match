package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	infra_postgres "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
)

type CompanyMemberRepo struct {
	db *sqlx.DB
}

func NewCompanyMemberRepo(db *sqlx.DB) *CompanyMemberRepo {
	return &CompanyMemberRepo{
		db: db,
	}
}

func (repo *CompanyMemberRepo) Create(ctx context.Context, member *domain.CompanyMember) error {
	exec := repo.getExec(ctx)
	query := `INSERT INTO company_members (user_id, company_id, role) VALUES ($1, $2, $3)`
	_, err := exec.ExecContext(ctx, query, member.UserID, member.CompanyID, member.Role)
	return err
}

// returns sqlx.TX if we're in transaction or r.db if not
func (repo *CompanyMemberRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value(infra_postgres.TxKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return repo.db
}
