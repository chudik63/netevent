package grpc

import (
	"context"
	"errors"
	"time"

	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/models"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/api/proto/event"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=eventservice.go -destination=mock/service_mock.go

type Service interface {
	CreateEvent(ctx context.Context, event *models.Event) (int64, error)
	ReadEvent(ctx context.Context, eventID int64) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	ListEvents(ctx context.Context) ([]*models.Event, error)
	ListEventsByCreator(ctx context.Context, creatorID int64) ([]*models.Event, error)
	CreateRegistration(ctx context.Context, userID int64, eventID int64) error
	SetChatStatus(ctx context.Context, participantID int64, eventID int64, isReady bool) error
	ListUsersToChat(ctx context.Context, eventID int64, userID int64) ([]*models.Participant, error)
	ListEventsByInterests(ctx context.Context, userID int64) ([]*models.Event, error)
	ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error)
	AddParticipant(ctx context.Context, participant *models.Participant) error
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

func (s *EventService) CreateEvent(ctx context.Context, req *event.CreateEventRequest) (*event.CreateEventResponse, error) {
	if _, err := time.Parse(models.TimeLayout, req.GetEvent().GetTime()); err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to create event", zap.String("err", models.ErrWrongTimeFormat.Error()))
		return nil, status.Errorf(codes.InvalidArgument, models.ErrWrongTimeFormat.Error())
	}

	resp, err := s.service.CreateEvent(ctx, &models.Event{
		CreatorID:   req.GetEvent().GetCreatorId(),
		Title:       req.GetEvent().GetTitle(),
		Description: req.GetEvent().GetDescription(),
		Time:        req.GetEvent().GetTime(),
		Place:       req.GetEvent().GetPlace(),
		Topics:      req.GetEvent().GetInterests(),
	})

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to create event", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.CreateEventResponse{
		EventId: resp,
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) DeleteEvent(ctx context.Context, req *event.DeleteEventRequest) (*event.DeleteEventResponse, error) {
	err := s.service.DeleteEvent(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to delete event", zap.String("err", err.Error()))
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.DeleteEventResponse{}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) ListEvents(ctx context.Context, req *event.ListEventsRequest) (*event.ListEventsResponse, error) {
	resp, err := s.service.ListEvents(ctx)

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.ListEventsResponse{
		Events: convertEventsToGRPC(resp),
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) ListEventsByCreator(ctx context.Context, req *event.ListEventsByCreatorRequest) (*event.ListEventsByCreatorResponse, error) {
	resp, err := s.service.ListEventsByCreator(ctx, req.GetCreatorId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.ListEventsByCreatorResponse{
		Events: convertEventsToGRPC(resp),
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) ListEventsByInterests(ctx context.Context, req *event.ListEventsByInterestsRequest) (*event.ListEventsByInterestsResponse, error) {
	resp, err := s.service.ListEventsByInterests(ctx, req.GetUserId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events by interests", zap.String("err", err.Error()))
		if errors.Is(err, models.ErrWrongUserId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.ListEventsByInterestsResponse{
		Events: convertEventsToGRPC(resp),
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) ListRegistratedEvents(ctx context.Context, req *event.ListRegistratedEventsRequest) (*event.ListRegistratedEventsResponse, error) {
	resp, err := s.service.ListRegistratedEvents(ctx, req.GetUserId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list events", zap.String("err", err.Error()))
		if errors.Is(err, models.ErrWrongUserId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.ListRegistratedEventsResponse{
		Events: convertEventsToGRPC(resp),
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) ListUsersToChat(ctx context.Context, req *event.ListUsersToChatRequest) (*event.ListUsersToChatResponse, error) {
	resp, err := s.service.ListUsersToChat(ctx, req.GetEventId(), req.GetUserId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to list users to chat", zap.String("err", err.Error()))
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
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
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) ReadEvent(ctx context.Context, req *event.ReadEventRequest) (*event.ReadEventResponse, error) {
	resp, err := s.service.ReadEvent(ctx, req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to read event", zap.String("err", err.Error()))
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
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
			Interests:   resp.Topics,
		},
	}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	err := s.service.CreateRegistration(ctx, req.GetUserId(), req.GetEventId())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to register user", zap.String("err", err.Error()))

		if errors.Is(err, models.ErrAlreadyRegistered) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, models.ErrWrongEventId) || errors.Is(err, models.ErrWrongUserId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.RegisterUserResponse{}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) (*event.UpdateEventResponse, error) {
	if _, err := time.Parse(models.TimeLayout, req.GetEvent().GetTime()); err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to create event", zap.String("err", models.ErrWrongTimeFormat.Error()))
		return nil, status.Errorf(codes.InvalidArgument, models.ErrWrongTimeFormat.Error())
	}

	err := s.service.UpdateEvent(ctx, &models.Event{
		EventID:     req.GetEvent().GetEventId(),
		CreatorID:   req.GetEvent().GetCreatorId(),
		Title:       req.GetEvent().GetTitle(),
		Description: req.GetEvent().GetDescription(),
		Time:        req.GetEvent().GetTime(),
		Place:       req.GetEvent().GetPlace(),
		Topics:      req.GetEvent().GetInterests(),
	})

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to update event", zap.String("err", err.Error()))
		if errors.Is(err, models.ErrWrongEventId) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.UpdateEventResponse{}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) SetChatStatus(ctx context.Context, req *event.SetChatStatusRequest) (*event.SetChatStatusResponse, error) {
	err := s.service.SetChatStatus(ctx, req.GetParticipantId(), req.GetEventId(), req.GetIsReady())

	if err != nil {
		s.logger.Error(context.WithValue(ctx, logger.RequestID, req.GetRequestId()), "failed to set chat status", zap.String("err", err.Error()))

		if errors.Is(err, models.ErrWrongEventId) || errors.Is(err, models.ErrWrongUserId) || errors.Is(err, models.ErrRegistrationNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.SetChatStatusResponse{}, status.New(codes.OK, "Success").Err()
}

func (s *EventService) AddParticipant(ctx context.Context, req *event.AddParticipantRequest) (*event.AddParticipantResponse, error) {
	err := s.service.AddParticipant(ctx, &models.Participant{
		UserID:    req.GetUser().GetUserId(),
		Name:      req.GetUser().GetName(),
		Email:     req.GetEmail(),
		Interests: req.GetUser().GetInterests(),
	})

	if err != nil {
		s.logger.Error(ctx, "failed to add participant", zap.String("err", err.Error()))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &event.AddParticipantResponse{}, status.New(codes.OK, "Success").Err()
}
