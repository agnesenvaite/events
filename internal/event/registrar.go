package event

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"

	"github.com/agnesenvaite/events/internal/event/classifier"
	"github.com/agnesenvaite/events/internal/event/invitee"
)

func NewCommandRegistrar(
	createHandler CreateHandler,
	createClassifierHandler classifier.CreateHandler,
	createInviteeHandler invitee.CreateHandler,
) error {
	if err := mediatr.RegisterRequestHandler[*CreateCommand, *Response](createHandler); err != nil {
		return errors.Wrap(err, "register create event command")
	}

	if err := mediatr.RegisterRequestHandler[*classifier.CreateCommand, []*classifier.Response](createClassifierHandler); err != nil {
		return errors.Wrap(err, "register create classifier command")
	}

	err := mediatr.RegisterRequestHandler[*invitee.CreateCommand, []*invitee.Response](createInviteeHandler)

	return errors.Wrap(err, "register create invitee command")
}
