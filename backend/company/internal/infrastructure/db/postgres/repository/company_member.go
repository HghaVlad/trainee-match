package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
)

type CompanyMemberRepo struct {
	db *sqlx.DB
}

func NewCompanyMemberRepo(db *sqlx.DB) *CompanyMemberRepo {
	return &CompanyMemberRepo{
		db: db,
	}
}

func (repo *CompanyMemberRepo) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	query := `SELECT * FROM company_members WHERE company_id = $1 AND user_id = $2`
	var member domain.CompanyMember
	err := repo.db.GetContext(ctx, &member, query, companyID, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain_errors.ErrCompanyMemberNotFound
	}

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (repo *CompanyMemberRepo) Create(ctx context.Context, member *domain.CompanyMember) error {
	exec := repo.getExec(ctx)
	query := `INSERT INTO company_members (user_id, company_id, role) VALUES ($1, $2, $3)`
	_, err := exec.ExecContext(ctx, query, member.UserID, member.CompanyID, member.Role)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain_errors.ErrCompanyMemberAlreadyExists
		case "23503":
			return domain_errors.ErrCompanyNotFound
		}
	}

	return err
}

func (repo *CompanyMemberRepo) UpdateRole(ctx context.Context, userID, companyID uuid.UUID, role value_types.CompanyRole) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`UPDATE company_members SET role = $1 WHERE user_id = $2 AND company_id = $3`,
		role, userID, companyID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return domain_errors.ErrCompanyMemberNotFound
	}

	return nil
}

func (repo *CompanyMemberRepo) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	exec := repo.getExec(ctx)

	res, err := exec.ExecContext(ctx,
		`DELETE FROM company_members WHERE user_id = $1 AND company_id = $2`,
		userID, companyID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return domain_errors.ErrCompanyMemberNotFound
	}

	return nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (repo *CompanyMemberRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := ctx.Value(infra_postgres.TxKey{}).(*sqlx.Tx)
	if ok {
		return tx
	}
	return repo.db
}
