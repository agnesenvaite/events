package api

import (
	"github.com/labstack/echo/v4"

	docsapi "github.com/agnesenvaite/events/internal/api/docs"
	eventapi "github.com/agnesenvaite/events/internal/api/event"
	statusapi "github.com/agnesenvaite/events/internal/api/status"
)

type RouteConfigurator func(*echo.Echo)

func NewRouteConfigurator(
	statusRouter *statusapi.Router,
	docsRouter *docsapi.Router,
	eventsRouter *eventapi.Router,
) RouteConfigurator {
	return func(e *echo.Echo) {
		v1 := e.Group("/v1")

		statusGroup := v1.Group("/status")
		statusRouter.Config(statusGroup)

		docsGroup := v1.Group("/docs")
		docsRouter.Config(docsGroup)

		eventsGroup := v1.Group("/events")
		eventsRouter.Config(eventsGroup)
	}
}
