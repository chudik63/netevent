package models

type Participant struct {
	UserID    int64
	Name      string
	Email     string
	Interests []string
}
