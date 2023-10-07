package invitee

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateCommand struct {
	Invitees []string
	EventID  string
}

type CreateHandler interface {
	Handle(ctx context.Context, command *CreateCommand) ([]*Response, error)
}

type createHandler struct {
	repository Repository
}

func NewCreateHandler(repository Repository) CreateHandler {
	return &createHandler{repository: repository}
}

func (h *createHandler) Handle(ctx context.Context, command *CreateCommand) ([]*Response, error) {
	var result []*Response

	now := time.Now().UTC()

	for i := range command.Invitees {
		invitee := &Invitee{
			ID:        uuid.NewString(),
			Invitee:   command.Invitees[i],
			EventID:   command.EventID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := h.repository.Create(ctx, invitee); err != nil {
			return nil, errors.Wrap(err, "create invitee")
		}

		result = append(result, invitee.toResponse())
	}

	return result, nil
}
