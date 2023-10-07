package container

import (
	_ "github.com/go-sql-driver/mysql" // driver import
	"go.uber.org/fx"

	"github.com/agnesenvaite/events/internal/api"
	docsapi "github.com/agnesenvaite/events/internal/api/docs"
	apierror "github.com/agnesenvaite/events/internal/api/error"
	eventapi "github.com/agnesenvaite/events/internal/api/event"
	statusapi "github.com/agnesenvaite/events/internal/api/status"
)

var APIModule = fx.Options(
	Module,
	fx.Provide(
		api.NewRouteConfigurator,
		apierror.NewErrorHandler,

		statusapi.NewController,
		statusapi.NewRouter,

		docsapi.NewRouter,

		eventapi.NewController,
		eventapi.NewRouter,
	),
)
