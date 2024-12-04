package service

import "event_service/internal/models"

type OrganizatorEventReposiory interface {
	CreateEvent(event models.Event) int64
	ReadEvent(eventID int64) models.Event
	UpdateEvent(event models.Event) error
	DeleteEvent(eventID int64) error
	ListEvents() []models.Event
}

type UserEventRepository interface {
	RegisterUser(participant models.Participant) error
	UpdateUser(participant models.Participant) error
	ListUsersToChat(eventID int64) []models.Participant
	ListEventsByUser(userID int64) []models.Event
}

type EventReposiory interface {
	OrganizatorEventReposiory
	UserEventRepository
}

type EventService struct {
	repository EventReposiory
}

func New(repo EventReposiory) *EventService {
	return &EventService{repository: repo}
}
