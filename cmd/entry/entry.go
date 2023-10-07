package entry

import (
	"context"

	"go.uber.org/fx"
)

type App interface {
	Options() []fx.Option
	Run(ctx context.Context) error
	Name() string
}
