package service

import (
	"context"
	"event_service/internal/models"
)

type OrganizatorEventReposiory interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context) []models.Event
}

type UserEventRepository interface {
	RegisterUser(ctx context.Context, participant *models.Participant) error
	UpdateUser(ctx context.Context, participant *models.Participant) error
	ListUsersToChat(ctx context.Context, eventID int64) []models.Participant
	ListEventsByUser(ctx context.Context, userID int64) []models.Event
}

type EventReposiory interface {
	OrganizatorEventReposiory
}

type EventService struct {
	repository EventReposiory
}

func New(repo EventReposiory) *EventService {
	return &EventService{repository: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, event *models.Event) (int64, error) {
	return 0, nil
}

func (s *EventService) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	return nil, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return nil
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID int64) error {
	return nil
}

func (s *EventService) ListEvents(ctx context.Context) []models.Event {
	return []models.Event{}
}

func (s *EventService) RegisterUser(ctx context.Context, participant *models.Participant) error {
	return nil
}

func (s *EventService) UpdateUser(ctx context.Context, participant *models.Participant) error {
	return nil
}

func (s *EventService) ListUsersToChat(ctx context.Context, eventID int64) []models.Participant {
	return []models.Participant{}
}

func (s *EventService) ListEventsByUser(ctx context.Context, userID int64) []models.Event {
	return []models.Event{}
}
