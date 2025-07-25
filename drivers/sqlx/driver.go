package sqlx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ctxKey string

const txKey ctxKey = "transactor_sqlx_tx"

var ErrNoTx = errors.New("no transaction in progress")

type Driver struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Driver {
	return &Driver{db: db}
}

func (*Driver) Name() string {
	return "sqlx"
}

func (d *Driver) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin: %w", err)
	}

	return d.setTx(ctx, tx), nil
}

func (d *Driver) Commit(ctx context.Context) error {
	tx, err := d.Tx(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (d *Driver) Rollback(ctx context.Context) error {
	tx, err := d.Tx(ctx)
	if err != nil {
		return err
	}

	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	return nil
}

func (*Driver) Tx(ctx context.Context) (*sqlx.Tx, error) {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx, nil
	}

	return nil, ErrNoTx
}

func (*Driver) setTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
