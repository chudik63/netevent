package models

import "time"

const TimeLayout = time.DateTime

type Event struct {
	EventID      int64
	CreatorID    int64
	Title        string
	Description  string
	Time         string
	Place        string
	Participants []Participant
	Topics       []string
}
