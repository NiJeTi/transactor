package transactor

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type Action func(ctx context.Context) error

type Transactor struct {
	drivers []Driver
}

func Init(drivers ...Driver) *Transactor {
	return &Transactor{
		drivers: drivers,
	}
}

func (t *Transactor) Do(ctx context.Context, action Action) error {
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
	for _, d := range t.drivers {
		dCtx, err := d.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to begin %s transaction: %w", d.Name(), err,
			)
		}

		ctx = dCtx
	}

	return ctx, nil
}

func (*Transactor) wrapAction(
	ctx context.Context, action Action,
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
	var resultErr error

	for _, d := range t.drivers {
		if err := t.wrapRollback(ctx, d); err != nil {
			resultErr = multierror.Append(resultErr, err)
		}
	}

	return resultErr
}

func (*Transactor) wrapRollback(ctx context.Context, d Driver) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%s rollback panic: %v", d.Name(), p)
		}
	}()

	if err := d.Rollback(ctx); err != nil {
		return fmt.Errorf("%s rollback error: %w", d.Name(), err)
	}

	return
}

func (t *Transactor) commit(ctx context.Context) error {
	var resultErr error

	for _, d := range t.drivers {
		if err := t.wrapCommit(ctx, d); err != nil {
			resultErr = multierror.Append(resultErr, err)
		}
	}

	return resultErr
}

func (*Transactor) wrapCommit(ctx context.Context, d Driver) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%s commit panic: %v", d.Name(), p)
		}
	}()

	if err := d.Commit(ctx); err != nil {
		return fmt.Errorf("%s commit error: %w", d.Name(), err)
	}

	return
}
