package classifier

import (
	"time"
)

var (
	LithuanianLanguage Language = "lithuanian"

	EnglishLanguage Language = "english"
	DutchLanguage   Language = "dutch"

	HDVideoQuality    VideoQuality = "720p"
	FullHDVideQuality VideoQuality = "1080p"
	UHDVideoQuality   VideoQuality = "2160p"

	HighAudioQuality   AudioQuality = "high"
	MediumAudioQuality AudioQuality = "medium"
	LowAudioQuality    AudioQuality = "low"

	LanguageType     Type = "language"
	VideoQualityType Type = "video_quality"
	AudioQualityType Type = "audio_quality"
)

type Language string
type VideoQuality string
type AudioQuality string
type Type string

type Classifier struct {
	ID        string    `db:"id"`
	Type      Type      `db:"type"`
	Value     string    `db:"value"`
	EventID   string    `db:"event_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Response struct {
	ID        string    `json:"id"`
	Type      Type      `json:"type"`
	Value     string    `json:"value"`
	EventID   string    `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (l Language) String() string {
	return string(l)
}

func (v VideoQuality) String() string {
	return string(v)
}

func (a AudioQuality) String() string {
	return string(a)
}

func (t Type) String() string {
	return string(t)
}

func (c *Classifier) toResponse() *Response {
	return &Response{
		ID:        c.ID,
		Type:      c.Type,
		Value:     c.Value,
		EventID:   c.EventID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
