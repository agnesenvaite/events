package event

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"

	_ "github.com/agnesenvaite/events/internal/api/error"
	"github.com/agnesenvaite/events/internal/config"
	"github.com/agnesenvaite/events/internal/event"
)

type Controller struct {
	cfg *config.Config
}

func NewController(cfg *config.Config) *Controller {
	return &Controller{
		cfg: cfg,
	}
}

// Create
// @Summary Create event
// @Description Create event
// @Tags Events
// @Produce json
// @Param data body createRequest true "Create request"
// @Success 200 {object} event.Response
// @Failure 400 {object} apierror.ListedError
// @Router /events [POST].
func (c *Controller) Create(ctx echo.Context) error {
	var request createRequest

	if err := ctx.Bind(&request); err != nil {
		return errors.Wrap(err, "bind request")
	}

	if err := request.validate(c.cfg.MaxInvitees); err != nil {
		return errors.Wrap(err, "validate")
	}

	result, err := mediatr.Send[*event.CreateCommand, *event.Response](ctx.Request().Context(), request.toCommand())
	if err != nil {
		return errors.Wrap(err, "send create command")
	}

	return ctx.JSON(http.StatusOK, result)
}
