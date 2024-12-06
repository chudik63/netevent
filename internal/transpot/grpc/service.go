package grpc

import (
	"context"
	"event_service/internal/logger"
	"event_service/internal/models"
	"event_service/pkg/api/proto/event"

	"go.uber.org/zap"
)

type Service interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context) []models.Event
	RegisterUser(ctx context.Context, participant *models.Participant) error
	UpdateUser(ctx context.Context, participant *models.Participant) error
	ListUsersToChat(ctx context.Context, eventID int64) []models.Participant
	ListEventsByUser(ctx context.Context, userID int64) []models.Event
}

type EventService struct {
	event.UnimplementedEventServiceServer
	service Service
	logger  logger.Logger
}

func NewEventService(ctx context.Context, s Service) *EventService {
	l := logger.GetLoggerFromCtx(ctx)

	return &EventService{
		service: s,
		logger:  l,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, req *event.CreateEventRequest) (*event.CreateEventResponse, error) {
	resp, err := s.service.CreateEvent(ctx, &models.Event{
		CreatorID:   req.Event.GetCreatorId(),
		Title:       req.Event.GetTitle(),
		Description: req.Event.GetDescription(),
		Time:        req.Event.GetTime(),
		Place:       req.Event.GetPlace(),
	})

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to create event", zap.String("err", err.Error()))
		return nil, err
	}

	return &event.CreateEventResponse{
		EventId: resp,
	}, nil
}

func (s *EventService) DeleteEvent(ctx context.Context, req *event.DeleteEventRequest) (*event.DeleteEventResponse, error) {
	return nil, nil
}

func (s *EventService) ListEvents(ctx context.Context, req *event.ListEventsRequest) (*event.ListEventsResponse, error) {
	return nil, nil
}

func (s *EventService) ListEventsByInterests(ctx context.Context, req *event.ListEventsByInterestsRequest) (*event.ListEventsByInterestsResponse, error) {
	return nil, nil
}

func (s *EventService) ListEventsByUser(ctx context.Context, req *event.ListEventsByUserRequest) (*event.ListEventsByUserResponse, error) {
	return nil, nil
}

func (s *EventService) ListUsersToChat(ctx context.Context, req *event.ListUsersToChatRequest) (*event.ListUsersToChatResponse, error) {
	return nil, nil
}

func (s *EventService) ReadEvent(ctx context.Context, req *event.ReadEventRequest) (*event.ReadEventResponse, error) {
	resp, err := s.service.ReadEvent(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to read event", zap.String("err", err.Error()))
		return nil, err
	}

	return &event.ReadEventResponse{
		Event: &event.Event{
			CreatorId:   resp.CreatorID,
			Title:       resp.Title,
			Description: resp.Description,
			Time:        resp.Time,
			Place:       resp.Place,
		},
	}, nil
}

func (s *EventService) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	return nil, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) (*event.UpdateEventResponse, error) {
	return nil, nil
}

func (s *EventService) UpdateUser(ctx context.Context, req *event.UpdateUserRequest) (*event.UpdateUserResponse, error) {
	return nil, nil
}
