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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/list"
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

func (repo *CompanyMemberRepo) ListViewsByCompany(
	ctx context.Context,
	companyID uuid.UUID,
	limit, offset int,
) ([]list.View, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT cm.user_id, cm.company_id, p.username, p.email, cm.role
		FROM company_members cm
		JOIN hr_user_projection p ON p.user_id = cm.user_id
		WHERE cm.company_id = $1
		ORDER BY username
		LIMIT $2 OFFSET $3`

	rows, err := q.Query(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	mems, err := scanMembersWithNames(rows)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	return mems, nil
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

func (repo *CompanyMemberRepo) GetCompanyRoleCount(
	ctx context.Context,
	companyID uuid.UUID,
	role member.CompanyRole,
) (int, error) {
	q := postgres.GetQuerier(ctx, repo.db)

	const query = `SELECT COUNT(user_id) 
		FROM company_members
		WHERE company_id = $1 AND role = $2`

	var cnt int

	err := q.QueryRow(ctx, query, companyID, role).Scan(&cnt)

	if err != nil {
		return 0, fmt.Errorf("get company member role count: %w", err)
	}

	return cnt, nil
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

func scanMembersWithNames(rows pgx.Rows) ([]list.View, error) {
	var mems []list.View

	for rows.Next() {
		var mem list.View

		err := rows.Scan(&mem.UserID, &mem.CompanyID, &mem.Username, &mem.Email, &mem.Role)
		if err != nil {
			return nil, fmt.Errorf("scan company member rows: %w", err)
		}

		mems = append(mems, mem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan company member rows: %w", err)
	}

	return mems, nil
}
