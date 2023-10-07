package status

import (
	"context"
	"database/sql"
	"time"
)

func NewMySQLChecker(db *sql.DB) Checker {
	return &MySQLChecker{db}
}

type MySQLChecker struct {
	db *sql.DB
}

func (mc *MySQLChecker) Name() string {
	return "mysql"
}

func (mc *MySQLChecker) Status(timeout time.Duration) string {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	if err := mc.db.PingContext(ctx); err != nil {
		return DOWN
	}

	return OK
}
