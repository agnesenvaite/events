package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/agnesenvaite/events/internal/transaction"
)

func Transactional(db *sqlx.DB, action func(tx transaction.DBTransaction) error) error {
	dbTransaction, err := db.Beginx()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	if err = action(dbTransaction); err != nil {
		if txErr := dbTransaction.Rollback(); txErr != nil {
			return errors.Wrap(txErr, "rollback transaction")
		}

		return err
	}

	return errors.Wrap(dbTransaction.Commit(), "commit transaction")
}
