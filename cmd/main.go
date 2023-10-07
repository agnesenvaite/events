package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/agnesenvaite/events/cmd/api"
	"github.com/agnesenvaite/events/cmd/entry"
	"github.com/agnesenvaite/events/internal/config"
	"github.com/agnesenvaite/events/internal/container"
)

type (
	ErrorChan *chan error
	StopChan  *chan os.Signal
)

type startupContext struct {
	Logger   *zap.Logger
	ErrChan  ErrorChan
	StopChan StopChan
}

const (
	APIContext = "api"
)

var Version = config.EnvLocal

func main() {
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Fatal("application context not provided")
	}

	start(applicationFactory())
}

func start(appContext entry.App) {
	var application startupContext

	app := fx.New(
		fx.NopLogger,
		fx.Provide(newStopChannel, newErrorChannel),
		fx.Options(appContext.Options()...),
		fx.Populate(&application.ErrChan, &application.StopChan, &application.Logger),
		fx.Supply(&container.AppSettings{Version: Version}),
	)

	if app.Err() != nil {
		log.Fatalf("%s failed to construct app\n %+v", appContext.Name(), errors.Wrap(app.Err(), "construct application"))
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	logger := application.Logger.With(zap.String("app", appContext.Name()))

	if err := app.Start(ctx); err != nil {
		logger.Fatal("failed to start application", zap.Error(errors.Wrap(err, "start app")))
	}

	// run application asynchronously and listen for errors channel, all fx OnStart lifecycle hooks are done by now
	go func() {
		logger.Info("starting application")

		*application.ErrChan <- errors.Wrap(appContext.Run(ctx), "run error")
	}()

	select {
	case err := <-*application.ErrChan:
		if err != nil {
			logger.Error("error while running", zap.Error(err))
		}
	case <-*application.StopChan:
		logger.Info("stopping application")
	}

	if err := app.Stop(ctx); err != nil {
		logger.Error("failed to stop app", zap.Error(errors.Wrap(err, "stop app")))
	}
}

func newStopChannel() StopChan {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	return &stop
}

func newErrorChannel() ErrorChan {
	errChan := make(chan error, 1)

	return &errChan
}

func applicationFactory() entry.App {
	var (
		applicationContext = flag.Arg(0)
		application        entry.App
	)

	switch applicationContext {
	case APIContext:
		application, _ = api.New()
	default:
		log.Fatalf("unknown run context '%s'", applicationContext)
	}

	return application
}
