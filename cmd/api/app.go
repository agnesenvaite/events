// @title Events API
// @version 1.0
// @description Events API specification

// Package api @host localhost
// @BasePath /v1
package api

import (
	"context"
	"fmt"

	echoLib "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	appapi "github.com/agnesenvaite/events/internal/api"
	apierror "github.com/agnesenvaite/events/internal/api/error"

	"github.com/agnesenvaite/events/cmd/entry"
	"github.com/agnesenvaite/events/internal/config"
	"github.com/agnesenvaite/events/internal/container"
)

type api struct {
	lc                fx.Lifecycle
	cfg               *config.Config
	logger            *zap.Logger
	routeConfigurator appapi.RouteConfigurator
	errorHandler      *apierror.Handler
}

func New() (entry.App, error) {
	return &api{}, nil
}

func (a *api) Options() []fx.Option {
	return []fx.Option{
		container.APIModule,
		fx.Invoke(a.setParams),
	}
}

func (a *api) Run(_ context.Context) error {
	engine := a.newEcho()

	a.lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			a.logger.Info("stopping echo engine")

			return errors.Wrap(engine.Shutdown(ctx), "stop engine")
		},
	})

	return errors.Wrap(engine.Start(fmt.Sprintf(":%d", +a.cfg.APIPort)), "run api engine")
}

func (a *api) Name() string {
	return "api"
}

func (a *api) setParams(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *zap.Logger,
	routeConfigurator appapi.RouteConfigurator,
	errorHandler *apierror.Handler,
) {
	a.lc = lc
	a.cfg = cfg
	a.logger = logger
	a.routeConfigurator = routeConfigurator
	a.errorHandler = errorHandler
}

func (a *api) newEcho() *echoLib.Echo {
	echoEngine := echoLib.New()
	echoEngine.HideBanner = true

	echoEngine.HTTPErrorHandler = a.errorHandler.Handler

	a.routeConfigurator(echoEngine)

	return echoEngine
}
