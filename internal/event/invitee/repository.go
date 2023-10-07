package invitee

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, invitee *Invitee) error
}
