package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/database/cache"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/kafka"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/models"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/repository"
)

//go:generate mockgen -source=event.go -destination=mock/mock.go

type Repository interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context, equations repository.Creds) ([]*models.Event, error)
	RegisterUser(ctx context.Context, userID int64, eventID int64) error
	SetChatStatus(ctx context.Context, userID int64, eventID int64, isReady bool) error
	ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error)
	ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error)
	InsertParticipant(ctx context.Context, participant *models.Participant) error
	ReadParticipant(ctx context.Context, userID int64) (*models.Participant, error)
	UpdateParticipant(ctx context.Context, participant *models.Participant) error
}

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Producer interface {
	Produce(ctx context.Context, message kafka.Message, topic string)
}

type EventService struct {
	repository Repository
	cache      Cache
	producer   Producer
}

func New(repo Repository, cache Cache, producer Producer) *EventService {
	return &EventService{
		repository: repo,
		cache:      cache,
		producer:   producer,
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
	_, err := s.repository.ReadParticipant(ctx, userID)
	if errors.Is(err, models.ErrWrongUserId) {
		return err
	}

	_, err = s.repository.ReadEvent(ctx, eventID)
	if errors.Is(err, models.ErrWrongEventId) {
		return err
	}

	err = s.repository.RegisterUser(ctx, userID, eventID)
	if err != nil {
		return err
	}

	event, err := s.repository.ReadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	participant, err := s.repository.ReadParticipant(ctx, userID)
	if err != nil {
		return err
	}

	s.producer.Produce(
		ctx,
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

func (s *EventService) SetChatStatus(ctx context.Context, userID int64, eventID int64, isReady bool) error {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if errors.Is(err, models.ErrWrongUserId) {
		return err
	}

	_, err = s.repository.ReadEvent(ctx, eventID)
	if errors.Is(err, models.ErrWrongEventId) {
		return err
	}

	return s.repository.SetChatStatus(ctx, userID, eventID, isReady)
}

func (s *EventService) ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error) {
	_, err := s.repository.ReadEvent(ctx, eventID)
	if errors.Is(err, models.ErrWrongEventId) {
		return nil, err
	}

	cacheKey := "readyToChat:" + strconv.FormatInt(eventID, 10)

	cachedData, err := s.cache.Get(ctx, cacheKey).Result()
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
		s.cache.Set(ctx, cacheKey, data, cache.Durability)
	}

	return participants, nil
}

func (s *EventService) ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error) {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if errors.Is(err, models.ErrWrongUserId) {
		return nil, err
	}

	return s.repository.ListEvents(ctx, repository.Creds{"user_id": userID})
}

func (s *EventService) ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error) {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if errors.Is(err, models.ErrWrongUserId) {
		return nil, err
	}

	return s.repository.ListEventsByInterests(ctx, userID)
}

func (s *EventService) AddParticipant(ctx context.Context, participant *models.Participant) error {
	_, err := s.repository.ReadParticipant(ctx, participant.UserID)
	if errors.Is(err, models.ErrWrongUserId) {
		_ = s.repository.InsertParticipant(ctx, participant)
		return nil
	}

	err = s.repository.UpdateParticipant(ctx, participant)

	return err
}
