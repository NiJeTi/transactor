package transactor

import (
	"context"
)

type Driver interface {
	Name() string

	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
