package transaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

type DBTransaction interface {
	sq.BaseRunner
	sq.ExecerContext

	Rollback() error
	Commit() error
}

type TransactionWrapper struct {
	transaction DBTransaction
}

type ExcludeMessage struct{}

type CommitHook func(context.Context) error

func (w *TransactionWrapper) Wrap(action func(transaction DBTransaction) error) error {
	return action(w.transaction)
}

type DBTransactionBehaviorHandler struct {
	db *sqlx.DB
}

func NewDBTransactionBehaviorHandler(db *sqlx.DB) *DBTransactionBehaviorHandler {
	return &DBTransactionBehaviorHandler{db: db}
}

func (h *DBTransactionBehaviorHandler) Handle(ctx context.Context, cmd any, next mediatr.RequestHandlerFunc) (any, error) {
	if _, ok := ctx.Value(TransactionWrapper{}).(*TransactionWrapper); ok {
		return next(ctx)
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "begin transaction")
	}

	ctx = context.WithValue(ctx, TransactionWrapper{}, &TransactionWrapper{transaction: tx})

	result, err := next(ctx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return result, errors.Wrap(rollbackErr, "rollback transaction")
		}

		return result, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return result, errors.Wrap(err, "commit transaction")
	}

	return result, nil
}
