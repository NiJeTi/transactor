package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type keyTx struct{}

var ErrNoTx = errors.New("no transaction in progress")

type Driver struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Driver {
	return &Driver{
		pool: pool,
	}
}

func (d *Driver) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin: %w", err)
	}

	return setTx(ctx, tx), nil
}

func (*Driver) Commit(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (*Driver) Rollback(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	if err := tx.Rollback(ctx); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	return nil
}

func (*Driver) Tx(ctx context.Context, def Tx) Tx {
	if tx, err := getTx(ctx); err == nil {
		return tx
	}

	return def
}

func getTx(ctx context.Context) (pgx.Tx, error) {
	if tx, ok := ctx.Value(keyTx{}).(pgx.Tx); ok {
		return tx, nil
	}

	return nil, ErrNoTx
}

func setTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, keyTx{}, tx)
}
