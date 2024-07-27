package sql

import (
	"context"
	"database/sql"
	"errors"
)

type rollbacker interface {
	Rollback() error
}

func RollbackOnPanicOrError(finalErr *error, rb rollbacker) {
	if finalErr == nil {
		panic(`must not pass nil final error to RollbackOnPanicOrError`)
	}

	if r := recover(); r != nil {
		*finalErr = errors.Join(rb.Rollback(), errors.New("recovered panic in RollbackOnPanicOrError"))
	} else if *finalErr != nil {
		*finalErr = errors.Join(rb.Rollback(), *finalErr)
	}
}

var (
	ErrNoRowsAffected       = errors.New("no rows affected")
	ErrMultipleRowsAffected = errors.New("multiple rows affected")
)

type execer interface {
	ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error)
}

// MustExecOne executes the given statement, then checks whether exactly one row was affected. If
// not exactly one row was affected, either ErrNoRowsAffected or ErrMultipleRowsAffected is
// returned.
func MustExecOne(
	ctx context.Context,
	execer execer,
	stmt string,
	args ...any,
) error {
	res, err := execer.ExecContext(ctx, stmt, args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	switch n {
	case 0:
		return ErrNoRowsAffected
	case 1:
		return nil
	default:
		return ErrMultipleRowsAffected
	}
}
