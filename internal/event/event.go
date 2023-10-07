package event

import (
	"time"

	"github.com/agnesenvaite/events/internal/event/classifier"
	"github.com/agnesenvaite/events/internal/event/invitee"
)

type Event struct {
	ID          string
	Name        string
	Date        time.Time
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Response struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	Languages      []string  `json:"languages"`
	VideoQualities []string  `json:"video_qualities"`
	AudioQualities []string  `json:"audio_qualities"`
	Invitees       []string  `json:"invitees"`
	Description    *string   `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (e *Event) toResponse(classifiers []*classifier.Response, invitees []*invitee.Response) *Response {
	response := &Response{
		ID:             e.ID,
		Name:           e.Name,
		Date:           e.Date,
		Languages:      []string{},
		VideoQualities: []string{},
		AudioQualities: []string{},
		Invitees:       []string{},
		Description:    e.Description,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}

	for i := range classifiers {
		switch classifiers[i].Type {
		case classifier.LanguageType:
			response.Languages = append(response.Languages, classifiers[i].Value)
		case classifier.VideoQualityType:
			response.VideoQualities = append(response.VideoQualities, classifiers[i].Value)
		case classifier.AudioQualityType:
			response.AudioQualities = append(response.AudioQualities, classifiers[i].Value)
		}
	}

	for i := range invitees {
		response.Invitees = append(response.Invitees, invitees[i].Invitee)
	}

	return response
}
