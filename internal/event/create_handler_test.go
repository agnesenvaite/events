package event

import (
	"context"
	"errors"
	"testing"
	"time"

	mock "github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"github.com/agnesenvaite/events/internal/config"
	"github.com/agnesenvaite/events/internal/event/classifier"
	"github.com/agnesenvaite/events/internal/event/invitee"
)

func TestHandle(t *testing.T) {
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockRepository(ctrl)
	mockConfig := &config.Config{
		DefaultVideoQuality: classifier.HDVideoQuality.String(),
		DefaultAudioQuality: classifier.HighAudioQuality.String(),
	}

	createClassifierHandler := classifier.NewMockCreateHandler(ctrl)
	createInviteeHandler := invitee.NewMockCreateHandler(ctrl)
	createHandler := NewCreateHandler(mockRepository, mockConfig)

	command := CreateCommand{
		Name:           "test name",
		Date:           time.Now().Add(time.Hour * 3),
		Languages:      []classifier.Language{classifier.DutchLanguage},
		VideoQualities: []classifier.VideoQuality{classifier.HDVideoQuality},
		AudioQualities: []classifier.AudioQuality{classifier.HighAudioQuality},
		Invitees:       []string{"invitee@test.com"},
		Description:    nil,
	}

	defaultAudioCommand := command
	defaultAudioCommand.AudioQualities = []classifier.AudioQuality{}
	defaultVideoCommand := command
	defaultVideoCommand.VideoQualities = []classifier.VideoQuality{}

	classifierResponse := []*classifier.Response{
		{
			Type:  classifier.LanguageType,
			Value: command.Languages[0].String(),
		},
		{
			Type:  classifier.VideoQualityType,
			Value: command.VideoQualities[0].String(),
		},
		{
			Type:  classifier.AudioQualityType,
			Value: command.AudioQualities[0].String(),
		},
	}

	inviteesResponse := []*invitee.Response{
		{
			Invitee: command.Invitees[0],
		},
	}

	testErr := errors.New("test err")

	err := NewCommandRegistrar(createHandler, createClassifierHandler, createInviteeHandler)
	assert.Nil(t, err)

	testCases := []struct {
		title          string
		command        *CreateCommand
		repositoryErr  error
		inviteesErr    error
		classifiersErr error
		handlerErr     error
	}{
		{
			title:   "happy path",
			command: &command,
		},
		{
			title:   "default audio",
			command: &defaultAudioCommand,
		},
		{
			title:   "default video",
			command: &defaultVideoCommand,
		},
		{
			title:         "repository error",
			command:       &command,
			repositoryErr: testErr,
			handlerErr:    testErr,
		},
		{
			title:       "invitees error",
			command:     &command,
			inviteesErr: testErr,
			handlerErr:  testErr,
		},
		{
			title:          "classifiers error",
			command:        &command,
			classifiersErr: testErr,
			handlerErr:     testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockRepository.EXPECT().Create(mock.Any(), mock.Any()).Return(tc.repositoryErr)

			if tc.repositoryErr == nil {
				createClassifierHandler.EXPECT().Handle(mock.Any(), mock.Any()).Return(classifierResponse, tc.classifiersErr)

				if tc.classifiersErr == nil {
					createInviteeHandler.EXPECT().Handle(mock.Any(), mock.Any()).Return(inviteesResponse, tc.inviteesErr)
				}
			}

			response, err := createHandler.Handle(context.Background(), tc.command)

			assert.ErrorIs(t, err, tc.handlerErr, tc.title)

			if err == nil {
				assert.Equal(t, command.Name, response.Name)
				assert.Equal(t, command.Languages[0].String(), response.Languages[0])
				assert.Equal(t, command.VideoQualities[0].String(), response.VideoQualities[0])
				assert.Equal(t, command.AudioQualities[0].String(), response.AudioQualities[0])
				assert.Equal(t, command.Invitees, response.Invitees)
				assert.Equal(t, command.Description, response.Description)
			}
		})
	}
}
