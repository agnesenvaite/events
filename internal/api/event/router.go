package event

import (
	echoLib "github.com/labstack/echo/v4"
)

type Router struct {
	controller *Controller
}

func NewRouter(controller *Controller) *Router {
	return &Router{controller: controller}
}

func (r *Router) Config(group *echoLib.Group) {
	group.POST("", r.controller.Create)
}
