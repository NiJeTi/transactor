package transactor

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type (
	keyNested struct{}
	valNested struct{}
)

type Transactor struct {
	driver Driver
}

func New(driver Driver) *Transactor {
	return &Transactor{
		driver: driver,
	}
}

func (t *Transactor) Do(
	ctx context.Context, action func(ctx context.Context) error,
) error {
	if isNested(ctx) {
		return action(ctx)
	}

	ctx = setIsNested(ctx)

	ctx, err := t.begin(ctx)
	if err != nil {
		return err
	}

	err = t.wrapAction(ctx, action)
	if err != nil {
		if rbErr := t.rollback(ctx); rbErr != nil {
			return multierror.Append(err, rbErr)
		}

		return err
	}

	err = t.commit(ctx)
	if err != nil {
		if rbErr := t.rollback(ctx); rbErr != nil {
			return multierror.Append(err, rbErr)
		}

		return err
	}

	return nil
}

func (t *Transactor) begin(ctx context.Context) (context.Context, error) {
	ctx, err := t.driver.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to begin transaction: %w", err,
		)
	}

	return ctx, nil
}

func (*Transactor) wrapAction(
	ctx context.Context, action func(ctx context.Context) error,
) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("action panic: %v", p)
		}
	}()

	if err := action(ctx); err != nil {
		return fmt.Errorf("action error: %w", err)
	}

	return
}

func (t *Transactor) rollback(ctx context.Context) error {
	if err := t.driver.Rollback(ctx); err != nil {
		return fmt.Errorf("rollback error: %w", err)
	}

	return nil
}

func (t *Transactor) commit(ctx context.Context) error {
	if err := t.driver.Commit(ctx); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

func isNested(ctx context.Context) bool {
	v := ctx.Value(keyNested{})

	return v != nil
}

func setIsNested(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyNested{}, valNested{})
}
