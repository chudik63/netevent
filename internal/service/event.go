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
	ListEventsByCreator(ctx context.Context, creatorID int64) ([]*models.Event, error)
}

type UserEventRepository interface {
	RegisterUser(ctx context.Context, participant *models.Participant) error
	UpdateUser(ctx context.Context, participant *models.Participant) error
	ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error)
	ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error)
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
	id, err := s.repository.CreateEvent(ctx, event)

	return id, err
}

func (s *EventService) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	event, err := s.repository.ReadEvent(ctx, eventID)

	return event, err
}

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return s.repository.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID int64) error {
	return s.repository.DeleteEvent(ctx, eventID)
}

func (s *EventService) ListEventsByCreator(ctx context.Context, creatorID int64) ([]*models.Event, error) {
	return s.repository.ListEventsByCreator(ctx, creatorID)
}

func (s *EventService) RegisterUser(ctx context.Context, participant *models.Participant) error {
	return nil
}

func (s *EventService) UpdateUser(ctx context.Context, participant *models.Participant) error {
	return nil
}

func (s *EventService) ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error) {
	return []*models.Participant{}, nil
}

func (s *EventService) ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error) {
	return []*models.Event{}, nil
}
