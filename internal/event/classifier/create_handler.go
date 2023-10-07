package classifier

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateCommand struct {
	Classifiers []*CreateCommandClassifier
	EventID     string
}

type CreateCommandClassifier struct {
	Type  Type
	Value string
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

	for i := range command.Classifiers {
		classifier := &Classifier{
			ID:        uuid.NewString(),
			Type:      command.Classifiers[i].Type,
			Value:     command.Classifiers[i].Value,
			EventID:   command.EventID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := h.repository.Create(ctx, classifier); err != nil {
			return nil, errors.Wrap(err, "create invitee")
		}

		result = append(result, classifier.toResponse())
	}

	return result, nil
}
