package models

import "time"

type Event struct {
	EventID      int64
	CreatorID    int64
	Title        string
	Description  string
	Time         time.Time
	Place        string
	Participants []Participant
}
