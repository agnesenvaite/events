package event

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/agnesenvaite/events/internal/event"
	"github.com/agnesenvaite/events/internal/event/classifier"
)

type createRequest struct {
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	Languages      []string  `json:"languages" enums:"english,lithuanian,dutch"`
	VideoQualities []string  `json:"video_qualities" enums:"720p,1080p,2160p"`
	AudioQualities []string  `json:"audio_qualities" enums:"high,medium,low"`
	Invitees       []string  `json:"invitees"`
	Description    *string   `json:"description,omitempty"`
}

func (r *createRequest) toCommand() *event.CreateCommand {
	var (
		languages      []classifier.Language
		videoQualities []classifier.VideoQuality
		audioQualities []classifier.AudioQuality
	)

	for i := range r.Languages {
		languages = append(languages, classifier.Language(r.Languages[i]))
	}

	for i := range r.VideoQualities {
		videoQualities = append(videoQualities, classifier.VideoQuality(r.VideoQualities[i]))
	}

	for i := range r.AudioQualities {
		audioQualities = append(audioQualities, classifier.AudioQuality(r.AudioQualities[i]))
	}

	return &event.CreateCommand{
		Name:           r.Name,
		Date:           r.Date,
		Languages:      languages,
		VideoQualities: videoQualities,
		AudioQualities: audioQualities,
		Invitees:       r.Invitees,
		Description:    r.Description,
	}
}

func (r *createRequest) validate(maxInvitees int) error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Date, validation.Required, validation.Min(time.Now().UTC())),
		validation.Field(
			&r.Languages,
			validation.Required,
			validation.Length(1, 100),
			validation.Each(
				validation.In(
					classifier.EnglishLanguage.String(),
					classifier.LithuanianLanguage.String(),
					classifier.DutchLanguage.String(),
				),
			),
		),
		validation.Field(
			&r.VideoQualities,
			validation.Length(0, 100),
			validation.Each(
				validation.In(
					classifier.HDVideoQuality.String(),
					classifier.FullHDVideQuality.String(),
					classifier.UHDVideoQuality.String(),
				),
			),
		),
		validation.Field(
			&r.AudioQualities,
			validation.Length(0, 100),
			validation.Each(
				validation.In(
					classifier.HighAudioQuality.String(),
					classifier.MediumAudioQuality.String(),
					classifier.LowAudioQuality.String(),
				),
			),
		),
		validation.Field(
			&r.Invitees,
			validation.Required,
			validation.Length(1, maxInvitees),
			validation.Each(is.Email, validation.Length(1, 100)),
		),
		validation.Field(&r.Description, validation.Length(1, 1000)),
	)
}
