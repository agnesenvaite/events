package container

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/blendle/zapdriver"
	_ "github.com/go-sql-driver/mysql" // driver import
	"github.com/jmoiron/sqlx"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/agnesenvaite/events/internal/config"
	"github.com/agnesenvaite/events/internal/event"
	"github.com/agnesenvaite/events/internal/event/classifier"
	"github.com/agnesenvaite/events/internal/event/invitee"
	"github.com/agnesenvaite/events/internal/status"
	"github.com/agnesenvaite/events/internal/transaction"
)

type AppSettings struct {
	Version string
}

var (
	defaultEnvFile = ".env"
	devEnvFile     = "dev.env"
	envFile        = flag.String("env-file", "", "Path to env file")
)

var Module = fx.Options(
	fx.Provide(
		newConfiguration,
		newDatabaseConnection,
		newWrappedDatabaseConnection,
		newLogger,
		newStatusPool,

		event.NewCreateHandler,
		event.NewMySQLRepository,

		classifier.NewCreateHandler,
		classifier.NewMySQLRepository,

		invitee.NewCreateHandler,
		invitee.NewMySQLRepository,

		transaction.NewDBTransactionBehaviorHandler,
	),
	fx.Invoke(
		newBehaviorRegistrar,
		event.NewCommandRegistrar,
	),
)

func newConfiguration() (*config.Config, error) {
	var cfg config.Config

	_, err := conf.Parse("", &cfg, config.FromEnvFiles(defaultEnvFile, devEnvFile, *envFile))

	return &cfg, err
}

func newDatabaseConnection(cfg *config.Config) (*sql.DB, error) {
	return sql.Open("mysql", cfg.Database.URL)
}

func newWrappedDatabaseConnection(conn *sql.DB, lc fx.Lifecycle) *sqlx.DB {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for i := 0; i < 4; i++ {
				if err := conn.Ping(); err == nil {
					return nil
				}

				select {
				case <-time.NewTicker(time.Second).C:
					continue
				case <-ctx.Done():
					return nil
				}
			}

			return errors.Wrap(conn.Ping(), "ping db connection on start")
		},
		OnStop: func(ctx context.Context) error {
			return errors.Wrap(conn.Close(), "close connection on stop")
		},
	})

	return sqlx.NewDb(conn, "mysql")
}

func newLogger(cfg *config.Config) (*zap.Logger, error) {
	errorLogs := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	infoLogs := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel
	})

	encoderConfig := zapdriver.NewDevelopmentEncoderConfig()

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	coreWrapper := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), infoLogs),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), errorLogs),
		)
	})

	return zapdriver.NewDevelopmentWithCore(coreWrapper)
}

func newStatusPool(db *sql.DB) status.Pool {
	return status.NewCheckerPool(
		status.NewMySQLChecker(db),
	)
}

func newBehaviorRegistrar(
	dbTransactionBehaviorHandler *transaction.DBTransactionBehaviorHandler,
) error {
	return errors.Wrap(
		mediatr.RegisterRequestPipelineBehaviors(
			dbTransactionBehaviorHandler,
		),
		"register pipeline behaviuors",
	)
}
