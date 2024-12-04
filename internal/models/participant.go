package models

type Participant struct {
	UserID    int64
	Name      string
	Interests []string
	Events    []Event
}
