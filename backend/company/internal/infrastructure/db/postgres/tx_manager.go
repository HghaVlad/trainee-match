package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool}
}

type txKeyT struct{}

//nolint:gochecknoglobals // ctx key
var txKey = txKeyT{}

func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// if already in tx - just execute
	if _, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	ctx = context.WithValue(ctx, txKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(context.WithoutCancel(ctx))
			panic(p)
		}
	}()

	err = fn(ctx)
	if err != nil {
		rollErr := tx.Rollback(context.WithoutCancel(ctx))
		if rollErr != nil {
			return errors.Join(err, rollErr)
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func GetQuerier(ctx context.Context, defaultQuerier Querier) Querier {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}

	return defaultQuerier
}
