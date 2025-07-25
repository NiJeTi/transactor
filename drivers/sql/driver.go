package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ctxKey string

const txKey ctxKey = "transactor_sql_tx"

var ErrNoTx = errors.New("no transaction in progress")

type Driver struct {
	db *sql.DB
}

func New(db *sql.DB) *Driver {
	return &Driver{db: db}
}

func (*Driver) Name() string {
	return "sql"
}

func (d *Driver) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.db.BeginTx(ctx, nil)
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

func (*Driver) Tx(ctx context.Context) (*sql.Tx, error) {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx, nil
	}

	return nil, ErrNoTx
}

func (*Driver) setTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
