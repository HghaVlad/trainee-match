package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/db/postgres"
)

type CompanyMemberRepo struct {
	db *pgxpool.Pool
}

func NewCompanyMemberRepo(db *pgxpool.Pool) *CompanyMemberRepo {
	return &CompanyMemberRepo{
		db: db,
	}
}

func (repo *CompanyMemberRepo) Get(
	ctx context.Context,
	userID, companyID uuid.UUID,
) (*member.CompanyMember, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT user_id, company_id, role	FROM company_members WHERE company_id = $1 AND user_id = $2`

	var memb member.CompanyMember

	err := q.QueryRow(ctx, query, companyID, userID).
		Scan(
			&memb.UserID,
			&memb.CompanyID,
			&memb.Role,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, member.ErrCompanyMemberNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get company member: %w", err)
	}

	return &memb, nil
}

func (repo *CompanyMemberRepo) Create(ctx context.Context, memb *member.CompanyMember) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `INSERT INTO company_members (user_id, company_id, role) VALUES ($1, $2, $3)`

	_, err := q.Exec(ctx, query, memb.UserID, memb.CompanyID, memb.Role)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return member.ErrCompanyMemberAlreadyExists
		case "23503":
			return company.ErrCompanyNotFound
		}
	}

	if err != nil {
		return fmt.Errorf("create company member: %w", err)
	}

	return nil
}

func (repo *CompanyMemberRepo) UpdateRole(
	ctx context.Context,
	userID, companyID uuid.UUID,
	role member.CompanyRole,
) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `UPDATE company_members SET role = $1 WHERE user_id = $2 AND company_id = $3`

	cmd, err := q.Exec(ctx, query, role, userID, companyID)

	if err != nil {
		return fmt.Errorf("update company member role: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return member.ErrCompanyMemberNotFound
	}

	return nil
}

func (repo *CompanyMemberRepo) Delete(
	ctx context.Context,
	userID, companyID uuid.UUID,
) error {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `DELETE FROM company_members WHERE user_id = $1 AND company_id = $2`

	cmd, err := q.Exec(ctx, query, userID, companyID)

	if err != nil {
		return fmt.Errorf("delete company member: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return member.ErrCompanyMemberNotFound
	}

	return nil
}
