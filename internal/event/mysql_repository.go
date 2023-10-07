package event

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/agnesenvaite/events/internal/db"
	"github.com/agnesenvaite/events/internal/transaction"
)

var (
	tableName = "events"

	colID          = "id"
	colName        = "name"
	colDate        = "date"
	colDescription = "description"
	colCreatedAt   = "created_at"
	colUpdatedAt   = "updated_at"

	columns = []string{
		colID,
		colName,
		colDate,
		colDescription,
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

func (r *mysqlRepository) Create(ctx context.Context, event *Event) error {
	values := []any{
		event.ID,
		event.Name,
		event.Date,
		event.Description,
		event.CreatedAt,
		event.UpdatedAt,
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
