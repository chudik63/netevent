package grpc

import (
	"context"
	"event_service/pkg/api/proto/event"
)

type Service interface {
}

type EventService struct {
	event.UnimplementedEventServiceServer
	service Service
}

func NewEventService(s Service) *EventService {
	return &EventService{service: s}
}

func (s *EventService) CreateEvent(context.Context, *event.CreateEventRequest) (*event.CreateEventResponse, error) {
	return nil, nil
}

func (s *EventService) DeleteEvent(context.Context, *event.DeleteEventRequest) (*event.DeleteEventResponse, error) {
	return nil, nil
}

func (s *EventService) ListEvents(context.Context, *event.ListEventsRequest) (*event.ListEventsResponse, error) {
	return nil, nil
}

func (s *EventService) ListEventsByInterests(context.Context, *event.ListEventsByInterestsRequest) (*event.ListEventsByInterestsResponse, error) {
	return nil, nil
}

func (s *EventService) ListEventsByUser(context.Context, *event.ListEventsByUserRequest) (*event.ListEventsByUserResponse, error) {
	return nil, nil
}

func (s *EventService) ListUsersToChat(context.Context, *event.ListUsersToChatRequest) (*event.ListUsersToChatResponse, error) {
	return nil, nil
}

func (s *EventService) ReadEvent(context.Context, *event.ReadEventRequest) (*event.ReadEventResponse, error) {
	return nil, nil
}

func (s *EventService) RegisterUser(context.Context, *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	return nil, nil
}

func (s *EventService) UpdateEvent(context.Context, *event.UpdateEventRequest) (*event.UpdateEventResponse, error) {
	return nil, nil
}

func (s *EventService) UpdateUser(context.Context, *event.UpdateUserRequest) (*event.UpdateUserResponse, error) {
	return nil, nil
}
