package progress

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -source transaction.go -destination transaction_mock.go -package progress

type transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}
