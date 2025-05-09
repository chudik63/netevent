package grpc

import (
	"context"
	"errors"

	"github.com/chudik63/netevent/events_service/internal/models"
	"github.com/chudik63/netevent/events_service/pkg/api/proto/event"
	"github.com/chudik63/netevent/events_service/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=eventservice.go -destination=mock/service_mock.go
type Service interface {
	CreateEvent(ctx context.Context, req *event.CreateEventRequest) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) error
	DeleteEvent(ctx context.Context, req *event.DeleteEventRequest) error
	ListEvents(ctx context.Context, creatorID int64) ([]*models.Event, error)
	ListEventsByInterests(ctx context.Context, req *event.ListEventsByInterestsRequest) ([]*models.Event, error)
	CreateRegistration(ctx context.Context, req *event.RegisterUserRequest) error
	SetChatStatus(ctx context.Context, req *event.SetChatStatusRequest) error
	ListUsersToChat(ctx context.Context, req *event.ListUsersToChatRequest) ([]*models.Participant, error)
	ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error)
	AddParticipant(ctx context.Context, req *event.AddParticipantRequest) error
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

func convertEventsToGRPC(events []*models.Event) []*event.Event {
	grpcEvents := make([]*event.Event, 0, len(events))
	for _, e := range events {
		grpcEvents = append(grpcEvents, &event.Event{
			EventId:     e.EventID,
			CreatorId:   e.CreatorID,
			Title:       e.Title,
			Description: e.Description,
			Time:        e.Time,
			Place:       e.Place,
			Interests:   e.Topics,
		})
	}
	return grpcEvents
}

func getRequestId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	i := md.Get("X-Request-ID")[0]

	return i
}

func (s *EventService) CreateEvent(ctx context.Context, req *event.CreateEventRequest) (*event.CreateEventResponse, error) {
	resp, err := s.service.CreateEvent(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrWrongTimeFormat) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to create event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.CreateEventResponse{
		EventId: resp,
	}, nil
}

func (s *EventService) DeleteEvent(ctx context.Context, req *event.DeleteEventRequest) (*event.DeleteEventResponse, error) {
	err := s.service.DeleteEvent(ctx, req)

	if err != nil {
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, models.ErrAccessDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to delete event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.DeleteEventResponse{}, nil
}

func (s *EventService) ListEvents(ctx context.Context, req *event.ListEventsRequest) (*event.ListEventsResponse, error) {
	resp, err := s.service.ListEvents(ctx, req.GetCreatorId())

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to list events", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.ListEventsResponse{
		Events: convertEventsToGRPC(resp),
	}, nil
}

func (s *EventService) ListEventsByInterests(ctx context.Context, req *event.ListEventsByInterestsRequest) (*event.ListEventsByInterestsResponse, error) {
	resp, err := s.service.ListEventsByInterests(ctx, req)

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to list events by interests", zap.Error(err))

		if errors.Is(err, models.ErrWrongUserId) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.ListEventsByInterestsResponse{
		Events: convertEventsToGRPC(resp),
	}, nil
}

func (s *EventService) ListRegistratedEvents(ctx context.Context, req *event.ListRegistratedEventsRequest) (*event.ListRegistratedEventsResponse, error) {
	resp, err := s.service.ListRegistratedEvents(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, models.ErrWrongUserId) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to list events", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.ListRegistratedEventsResponse{
		Events: convertEventsToGRPC(resp),
	}, nil
}

func (s *EventService) ListUsersToChat(ctx context.Context, req *event.ListUsersToChatRequest) (*event.ListUsersToChatResponse, error) {
	resp, err := s.service.ListUsersToChat(ctx, req)

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to list users to chat", zap.Error(err))
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	participants := make([]*event.Participant, 0, len(resp))
	for _, p := range resp {
		participants = append(participants, &event.Participant{
			UserId:    p.UserID,
			Name:      p.Name,
			Interests: p.Interests,
		})
	}

	return &event.ListUsersToChatResponse{
		Participants: participants,
	}, nil
}

func (s *EventService) ReadEvent(ctx context.Context, req *event.ReadEventRequest) (*event.ReadEventResponse, error) {
	resp, err := s.service.ReadEvent(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to read event", zap.Error(err))
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.ReadEventResponse{
		Event: &event.Event{
			EventId:     resp.EventID,
			CreatorId:   resp.CreatorID,
			Title:       resp.Title,
			Description: resp.Description,
			Time:        resp.Time,
			Place:       resp.Place,
			Interests:   resp.Topics,
		},
	}, nil
}

func (s *EventService) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	err := s.service.CreateRegistration(ctx, req)

	if err != nil {
		if errors.Is(err, models.ErrAlreadyRegistered) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, models.ErrWrongEventId) || errors.Is(err, models.ErrWrongUserId) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to register user", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.RegisterUserResponse{}, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) (*event.UpdateEventResponse, error) {
	err := s.service.UpdateEvent(ctx, req)

	if err != nil {
		if errors.Is(err, models.ErrWrongTimeFormat) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, models.ErrAccessDenied) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to update event", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.UpdateEventResponse{}, nil
}

func (s *EventService) SetChatStatus(ctx context.Context, req *event.SetChatStatusRequest) (*event.SetChatStatusResponse, error) {
	err := s.service.SetChatStatus(ctx, req)

	if err != nil {
		if errors.Is(err, models.ErrWrongEventId) || errors.Is(err, models.ErrWrongUserId) || errors.Is(err, models.ErrRegistrationNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		s.logger.Error(context.WithValue(ctx, logger.RequestID, getRequestId(ctx)), "failed to set chat status", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.SetChatStatusResponse{}, nil
}

func (s *EventService) AddParticipant(ctx context.Context, req *event.AddParticipantRequest) (*event.AddParticipantResponse, error) {
	err := s.service.AddParticipant(ctx, req)

	if err != nil {
		s.logger.Error(ctx, "failed to add participant", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &event.AddParticipantResponse{}, nil
}
