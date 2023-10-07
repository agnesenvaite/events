package classifier

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/agnesenvaite/events/internal/db"
	"github.com/agnesenvaite/events/internal/transaction"
)

var (
	tableName = "event_classifiers"

	colID        = "id"
	colType      = "type"
	colValue     = "value"
	colEventID   = "event_id"
	colCreatedAt = "created_at"
	colUpdatedAt = "updated_at"

	columns = []string{
		colID,
		colType,
		colValue,
		colEventID,
		colCreatedAt,
		colUpdatedAt,
	}
)

type mysqlRepository struct {
	mysqlDB *sqlx.DB
}

func NewMySQLRepository(mysqlDB *sqlx.DB) Repository {
	return &mysqlRepository{mysqlDB: mysqlDB}
}

func (r *mysqlRepository) Create(ctx context.Context, classifier *Classifier) error {
	values := []any{
		classifier.ID,
		classifier.Type,
		classifier.Value,
		classifier.EventID,
		classifier.CreatedAt,
		classifier.UpdatedAt,
	}

	action := func(dbTransaction transaction.DBTransaction) error {
		builder := sq.Insert(tableName).Columns(columns...).Values(values...)

		if transactionCtx, ok := ctx.Value(transaction.TransactionWrapper{}).(*transaction.TransactionWrapper); ok {
			return transactionCtx.Wrap(func(dbTransaction transaction.DBTransaction) error {
				_, err := builder.RunWith(dbTransaction).ExecContext(ctx)

				return errors.Wrap(err, "execute query")
			})
		}

		_, err := builder.RunWith(r.mysqlDB).ExecContext(ctx)

		return errors.Wrap(err, "execute query")
	}

	if transactionCtx, ok := ctx.Value(transaction.TransactionWrapper{}).(*transaction.TransactionWrapper); ok {
		return transactionCtx.Wrap(action)
	}

	return errors.Wrap(db.Transactional(r.mysqlDB, action), "execute in transaction")
}
