package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/chudik63/netevent/events_service/internal/database/cache"
	"github.com/chudik63/netevent/events_service/internal/models"
	"github.com/chudik63/netevent/events_service/internal/producer"
	"github.com/chudik63/netevent/events_service/internal/repository"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=event.go -destination=mock/mock.go

type Repository interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context, equations repository.Creds) ([]*models.Event, error)
	ListEventsByInterests(ctx context.Context, userID int64, equations repository.Creds) ([]*models.Event, error)
	CreateRegistration(ctx context.Context, userID int64, eventID int64) error
	ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error)
	SetChatStatus(ctx context.Context, userID int64, eventID int64, isReady bool) error
	ListUsersToChat(ctx context.Context, eventID int64, userID int64) ([]*models.Participant, error)
	CreateParticipant(ctx context.Context, participant *models.Participant) error
	ReadParticipant(ctx context.Context, userID int64) (*models.Participant, error)
	UpdateParticipant(ctx context.Context, participant *models.Participant) error
}

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Producer interface {
	Produce(ctx context.Context, message producer.Message, topic string)
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

func (s *EventService) UpdateEvent(ctx context.Context, event *models.Event, userID int64) error {
	if userID != event.CreatorID {
		return models.ErrAccessDenied
	}

	return s.repository.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID int64, userID int64) error {
	event, err := s.repository.ReadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	if userID != event.CreatorID {
		return models.ErrAccessDenied
	}

	return s.repository.DeleteEvent(ctx, eventID)
}

func (s *EventService) ListEvents(ctx context.Context, creds repository.Creds) ([]*models.Event, error) {
	return s.repository.ListEvents(ctx, creds)
}

func (s *EventService) CreateRegistration(ctx context.Context, userID int64, eventID int64) error {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if err != nil {
		return err
	}

	_, err = s.repository.ReadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	err = s.repository.CreateRegistration(ctx, userID, eventID)
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
		producer.Message{
			UserEmail:  participant.Email,
			UserName:   participant.Name,
			EventName:  event.Title,
			EventTime:  event.Time,
			EventPlace: event.Place,
		},
		producer.RegistrationTopic,
	)

	return nil
}

func (s *EventService) SetChatStatus(ctx context.Context, userID int64, eventID int64, isReady bool) error {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if err != nil {
		return err
	}

	_, err = s.repository.ReadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	err = s.repository.CreateRegistration(ctx, userID, eventID)
	if errors.Is(err, models.ErrAlreadyRegistered) {
		return s.repository.SetChatStatus(ctx, userID, eventID, isReady)
	}

	return models.ErrRegistrationNotFound
}

func (s *EventService) ListUsersToChat(ctx context.Context, eventID int64, userID int64) ([]*models.Participant, error) {
	_, err := s.repository.ReadEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	participants, err := s.repository.ListUsersToChat(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (s *EventService) ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error) {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if errors.Is(err, models.ErrWrongUserId) {
		return nil, err
	}

	return s.repository.ListRegistratedEvents(ctx, userID)
}

func (s *EventService) ListEventsByInterests(ctx context.Context, userID int64, creds repository.Creds) ([]*models.Event, error) {
	_, err := s.repository.ReadParticipant(ctx, userID)
	if errors.Is(err, models.ErrWrongUserId) {
		return nil, err
	}

	cacheKey := strconv.FormatInt(userID, 10) + "_eventsByInterests"
	for cred, value := range creds {
		cacheKey += fmt.Sprintf("_%s=%s", cred, value)
	}

	cachedData, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var events []*models.Event
		if err := json.Unmarshal([]byte(cachedData), &events); err == nil {
			return events, nil
		}
	}

	events, err := s.repository.ListEventsByInterests(ctx, userID, creds)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(events)
	if err == nil {
		s.cache.Set(ctx, cacheKey, data, cache.Durability)
	}

	return events, nil
}

func (s *EventService) AddParticipant(ctx context.Context, participant *models.Participant) error {
	_, err := s.repository.ReadParticipant(ctx, participant.UserID)
	if errors.Is(err, models.ErrWrongUserId) {
		return s.repository.CreateParticipant(ctx, participant)
	}

	err = s.repository.UpdateParticipant(ctx, participant)

	return err
}
