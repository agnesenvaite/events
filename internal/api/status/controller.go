package status

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/agnesenvaite/events/internal/status"
)

type Controller struct {
	checkerPool status.Pool
}

func NewController(pool status.Pool) *Controller {
	return &Controller{checkerPool: pool}
}

// Ping godoc
// @Summary Get service status
// @Description Check if service is in running state and can access its resources
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} status.ResponseStatus
// @Router /ping [GET].
func (ec *Controller) Ping(context echo.Context) error {
	currentStatus := ec.checkerPool.Status()
	if currentStatus == status.DOWN {
		return context.JSON(http.StatusServiceUnavailable, ResponseStatus{Status: ec.checkerPool.Status()})
	}

	return context.JSON(http.StatusOK, ResponseStatus{Status: ec.checkerPool.Status()})
}

// Details godoc
// @Summary Get service detailed status
// @Description Check if service is in running state and can access its resources with additional list of resource and its status
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} status.ResponseDetails
// @Router /details [GET].
func (ec *Controller) Details(context echo.Context) error {
	return context.JSON(http.StatusOK, ResponseDetails{Status: ec.checkerPool.Details()})
}
