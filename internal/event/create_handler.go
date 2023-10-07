package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"

	"github.com/agnesenvaite/events/internal/config"
	"github.com/agnesenvaite/events/internal/event/classifier"
	"github.com/agnesenvaite/events/internal/event/invitee"
)

type CreateCommand struct {
	Name           string
	Date           time.Time
	Languages      []classifier.Language
	VideoQualities []classifier.VideoQuality
	AudioQualities []classifier.AudioQuality
	Invitees       []string
	Description    *string
}

type CreateHandler interface {
	Handle(ctx context.Context, command *CreateCommand) (*Response, error)
}

type createHandler struct {
	repository Repository
	cfg        *config.Config
}

func NewCreateHandler(repository Repository, cfg *config.Config) CreateHandler {
	return &createHandler{
		repository: repository,
		cfg:        cfg,
	}
}

func (h *createHandler) Handle(ctx context.Context, command *CreateCommand) (*Response, error) {
	now := time.Now().UTC()

	event := &Event{
		ID:          uuid.NewString(),
		Name:        command.Name,
		Date:        command.Date.UTC(),
		Description: command.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if len(command.AudioQualities) == 0 {
		command.AudioQualities = []classifier.AudioQuality{classifier.AudioQuality(h.cfg.DefaultAudioQuality)}
	}

	if len(command.VideoQualities) == 0 {
		command.VideoQualities = []classifier.VideoQuality{classifier.VideoQuality(h.cfg.DefaultVideoQuality)}
	}

	if err := h.repository.Create(ctx, event); err != nil {
		return nil, errors.Wrap(err, "create in repo")
	}

	classifiers, err := h.createClassifiers(ctx, event.ID, command)
	if err != nil {
		return nil, errors.Wrap(err, "create classifiers")
	}

	invitees, err := h.createInvitees(ctx, event.ID, command)
	if err != nil {
		return nil, errors.Wrap(err, "crete invitees")
	}

	return event.toResponse(classifiers, invitees), nil
}

func (h *createHandler) createClassifiers(ctx context.Context, eventID string, command *CreateCommand) ([]*classifier.Response, error) {
	classifierCommand := &classifier.CreateCommand{
		Classifiers: []*classifier.CreateCommandClassifier{},
		EventID:     eventID,
	}

	for i := range command.Languages {
		classifierCommand.Classifiers = append(classifierCommand.Classifiers, &classifier.CreateCommandClassifier{
			Type:  classifier.LanguageType,
			Value: command.Languages[i].String(),
		})
	}

	for i := range command.VideoQualities {
		classifierCommand.Classifiers = append(classifierCommand.Classifiers, &classifier.CreateCommandClassifier{
			Type:  classifier.VideoQualityType,
			Value: command.VideoQualities[i].String(),
		})
	}

	for i := range command.AudioQualities {
		classifierCommand.Classifiers = append(classifierCommand.Classifiers, &classifier.CreateCommandClassifier{
			Type:  classifier.AudioQualityType,
			Value: command.AudioQualities[i].String(),
		})
	}

	result, err := mediatr.Send[*classifier.CreateCommand, []*classifier.Response](ctx, classifierCommand)

	return result, errors.Wrap(err, "create classifiers")
}

func (h *createHandler) createInvitees(ctx context.Context, eventID string, command *CreateCommand) ([]*invitee.Response, error) {
	inviteeCommand := &invitee.CreateCommand{
		Invitees: command.Invitees,
		EventID:  eventID,
	}

	result, err := mediatr.Send[*invitee.CreateCommand, []*invitee.Response](ctx, inviteeCommand)

	return result, errors.Wrap(err, "create invitees")
}
