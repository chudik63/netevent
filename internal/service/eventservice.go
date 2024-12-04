package service

type EventReposiory interface {
}

type EventService struct {
	repository EventReposiory
}

func New(repo EventReposiory) *EventService {
	return &EventService{repository: repo}
}
