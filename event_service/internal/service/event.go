package service

import (
	"context"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/models"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/repository"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context, equations repository.Creds) ([]*models.Event, error)
	RegisterUser(ctx context.Context, userID int64, eventID int64) error
	SetChatStatus(ctx context.Context, participantID int64, eventID int64, isReady bool) error
	ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error)
	ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error)
}

type EventService struct {
	repository Repository
}

func New(repo Repository) *EventService {
	return &EventService{repository: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, event *models.Event) (int64, error) {
	return s.repository.CreateEvent(ctx, event)
}

func (s *EventService) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	return s.repository.ReadEvent(ctx, eventID)
}

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return s.repository.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID int64) error {
	return s.repository.DeleteEvent(ctx, eventID)
}

func (s *EventService) ListEvents(ctx context.Context) ([]*models.Event, error) {
	return s.repository.ListEvents(ctx, repository.Creds{})
}

func (s *EventService) ListEventsByCreator(ctx context.Context, creatorID int64) ([]*models.Event, error) {
	return s.repository.ListEvents(ctx, repository.Creds{"creator_id": creatorID})
}

func (s *EventService) RegisterUser(ctx context.Context, userID int64, eventID int64) error {
	return s.repository.RegisterUser(ctx, userID, eventID)
}

func (s *EventService) SetChatStatus(ctx context.Context, participantID int64, eventID int64, isReady bool) error {
	return s.repository.SetChatStatus(ctx, participantID, eventID, isReady)
}

func (s *EventService) ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error) {
	return s.repository.ListUsersToChat(ctx, eventID)
}

func (s *EventService) ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error) {
	return s.repository.ListEvents(ctx, repository.Creds{"user_id": userID})
}

func (s *EventService) ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error) {
	return s.repository.ListEventsByInterests(ctx, userID)
}
