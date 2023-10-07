package invitee

import (
	"time"
)

type Invitee struct {
	ID        string    `db:"id"`
	Invitee   string    `db:"invitee"`
	EventID   string    `db:"event_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Response struct {
	ID        string    `json:"id"`
	Invitee   string    `json:"invitee"`
	EventID   string    `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *Invitee) toResponse() *Response {
	return &Response{
		ID:        i.ID,
		Invitee:   i.Invitee,
		EventID:   i.EventID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}
