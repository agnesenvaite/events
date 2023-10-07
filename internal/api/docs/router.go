package docs

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/agnesenvaite/events/docs"
	"github.com/agnesenvaite/events/internal/config"
)

type Router struct {
}

func NewRouter(cfg *config.Config) *Router {
	docs.SwaggerInfo.Host = cfg.Docs.Host

	return &Router{}
}

func (r *Router) Config(engine *echo.Group) {
	engine.GET("*", echoSwagger.WrapHandler)
}
