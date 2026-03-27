package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TxKey struct{}

type TxManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	tx, err := m.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	// rollback by default - if commit() - no changes
	defer func() {
		_ = tx.Rollback()
	}()

	ctx = context.WithValue(ctx, TxKey{}, tx)

	err = fn(ctx)

	if err != nil {
		return err
	}

	return tx.Commit()
}
