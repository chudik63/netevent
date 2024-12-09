package grpc

import (
	"context"
	"event_service/internal/logger"
	"event_service/internal/models"
	"event_service/pkg/api/proto/event"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context) ([]*models.Event, error)
	ListEventsByCreator(ctx context.Context, creatorID int64) ([]*models.Event, error)
	RegisterUser(ctx context.Context, participant *models.Participant, eventID int64) error
	SetChatStatus(ctx context.Context, participantID int64, eventID int64, isReady bool) error
	ListUsersToChat(ctx context.Context, eventID int64) ([]*models.Participant, error)
	ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error)
	ListEventsByUser(ctx context.Context, userID int64) ([]*models.Event, error)
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
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.CreateEventResponse{
		EventId: resp,
	}, nil
}

func (s *EventService) DeleteEvent(ctx context.Context, req *event.DeleteEventRequest) (*event.DeleteEventResponse, error) {
	err := s.service.DeleteEvent(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to delete event", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.DeleteEventResponse{
		Message: "OK",
	}, nil
}

func (s *EventService) ListEvents(ctx context.Context, req *event.ListEventsRequest) (*event.ListEventsResponse, error) {
	resp, err := s.service.ListEvents(ctx)

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	events := make([]*event.Event, 0, len(resp))
	for _, e := range resp {
		events = append(events, &event.Event{
			EventId:     e.EventID,
			CreatorId:   e.CreatorID,
			Title:       e.Title,
			Description: e.Description,
			Time:        e.Time,
			Place:       e.Place,
		})
	}

	return &event.ListEventsResponse{
		Events: events,
	}, nil
}

func (s *EventService) ListEventsByCreator(ctx context.Context, req *event.ListEventsByCreatorRequest) (*event.ListEventsByCreatorResponse, error) {
	resp, err := s.service.ListEventsByCreator(ctx, req.GetCreatorId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	events := make([]*event.Event, 0, len(resp))
	for _, e := range resp {
		events = append(events, &event.Event{
			EventId:     e.EventID,
			CreatorId:   e.CreatorID,
			Title:       e.Title,
			Description: e.Description,
			Time:        e.Time,
			Place:       e.Place,
		})
	}

	return &event.ListEventsByCreatorResponse{
		Events: events,
	}, nil
}

func (s *EventService) ListEventsByInterests(ctx context.Context, req *event.ListEventsByInterestsRequest) (*event.ListEventsByInterestsResponse, error) {
	return nil, nil
}

func (s *EventService) ListEventsByUser(ctx context.Context, req *event.ListEventsByUserRequest) (*event.ListEventsByUserResponse, error) {
	resp, err := s.service.ListEventsByUser(ctx, req.GetUserId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	events := make([]*event.Event, 0, len(resp))
	for _, e := range resp {
		events = append(events, &event.Event{
			EventId:     e.EventID,
			CreatorId:   e.CreatorID,
			Title:       e.Title,
			Description: e.Description,
			Time:        e.Time,
			Place:       e.Place,
		})
	}

	return &event.ListEventsByUserResponse{
		Events: events,
	}, nil
}

func (s *EventService) ListUsersToChat(ctx context.Context, req *event.ListUsersToChatRequest) (*event.ListUsersToChatResponse, error) {
	resp, err := s.service.ListUsersToChat(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list users to chat", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	participants := make([]*event.Participant, 0, len(resp))
	for _, p := range resp {
		participants = append(participants, &event.Participant{
			UserId:    p.UserID,
			Name:      p.Name,
			Interests: strings.Join(p.Interests, ", "),
		})
	}

	return &event.ListUsersToChatResponse{
		Participants: participants,
	}, nil
}

func (s *EventService) ReadEvent(ctx context.Context, req *event.ReadEventRequest) (*event.ReadEventResponse, error) {
	resp, err := s.service.ReadEvent(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to read event", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.ReadEventResponse{
		Event: &event.Event{
			EventId:     resp.EventID,
			CreatorId:   resp.CreatorID,
			Title:       resp.Title,
			Description: resp.Description,
			Time:        resp.Time,
			Place:       resp.Place,
		},
	}, nil
}

func (s *EventService) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	err := s.service.RegisterUser(ctx, &models.Participant{
		UserID: req.GetParticipant().GetUserId(),
		Name:   req.GetParticipant().GetName(),
	},
		req.GetEventId(),
	)

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to register user", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.RegisterUserResponse{
		Message: "OK",
	}, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) (*event.UpdateEventResponse, error) {
	err := s.service.UpdateEvent(ctx, &models.Event{
		EventID:     req.Event.GetEventId(),
		CreatorID:   req.Event.GetCreatorId(),
		Title:       req.Event.GetTitle(),
		Description: req.Event.GetDescription(),
		Time:        req.Event.GetTime(),
		Place:       req.Event.GetPlace(),
	})

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to update event", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.UpdateEventResponse{
		Message: "OK",
	}, nil
}

func (s *EventService) SetChatStatus(ctx context.Context, req *event.SetChatStatusRequest) (*event.SetChatStatusResponse, error) {
	err := s.service.SetChatStatus(ctx, req.GetParticipantId(), req.GetEventId(), req.GetIsReady())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to set chat status", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.SetChatStatusResponse{
		Message: "OK",
	}, nil
}
