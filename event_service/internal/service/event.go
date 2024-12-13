package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/cache"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/kafka"
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
	ReadParticipant(ctx context.Context, userID int64) (*models.Participant, error)
}

type EventService struct {
	repository  Repository
	redisClient *redis.Client
	producer    *kafka.Producer
}

func New(repo Repository, redis *redis.Client, producer *kafka.Producer) *EventService {
	return &EventService{
		repository:  repo,
		redisClient: redis,
		producer:    producer,
	}
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
	err := s.repository.RegisterUser(ctx, userID, eventID)
	if err != nil {
		return err
	}

	event, err := s.repository.ReadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	participant, err := s.repository.ReadParticipant(ctx, userID)

	err = s.producer.Produce(
		kafka.Message{
			UserEmail:  participant.Email,
			UserName:   participant.Name,
			EventName:  event.Title,
			EventTime:  event.Time,
			EventPlace: event.Place,
		},
		kafka.RegistrationTopic,
	)

	return nil
}

func (s *EventService) SetChatStatus(ctx context.Context, participantID int64, eventID int64, isReady bool) error {
	return s.repository.SetChatStatus(ctx, participantID, eventID, isReady)
}

func (s *EventService) ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error) {
	cacheKey := "readyToChat:" + strconv.FormatInt(eventID, 10)

	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var participants []*models.Participant
		if err := json.Unmarshal([]byte(cachedData), &participants); err == nil {
			return participants, nil
		}
	}

	participants, err := s.repository.ListUsersToChat(ctx, eventID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(participants)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, data, cache.Durability)
	}

	return participants, nil
}

func (s *EventService) ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error) {
	return s.repository.ListEvents(ctx, repository.Creds{"user_id": userID})
}

func (s *EventService) ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error) {
	return s.repository.ListEventsByInterests(ctx, userID)
}
