package classifier

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, classifier *Classifier) error
}
