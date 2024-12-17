package grpc

import (
	"context"

	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/internal/client"
	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/pkg/api/gateway"
	auth "gitlab.crja72.ru/gospec/go9/netevent/auth/pkg/proto"
	event "gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/api/proto/event"
)

type GatewayServer struct {
	gateway.UnimplementedGatewayServer
	authClient  *client.AuthClient
	eventClient *client.EventClient
}

func convertEvents(s []*event.Event) []*gateway.Event {
	events := make([]*gateway.Event, len(s))

	for i, e := range s {
		events[i] = &gateway.Event{
			EventId:     e.GetEventId(),
			CreatorId:   e.GetCreatorId(),
			Title:       e.GetTitle(),
			Description: e.GetDescription(),
			Time:        e.GetTime(),
			Place:       e.GetPlace(),
			Interests:   e.GetInterests(),
		}
	}

	return events
}

func convertUsers(s []*event.Participant) []*gateway.Participant {
	users := make([]*gateway.Participant, len(s))

	for i, p := range s {
		users[i] = &gateway.Participant{
			UserId:    p.GetUserId(),
			Name:      p.GetName(),
			Interests: p.GetInterests(),
		}
	}

	return users
}

func NewGatewayServer(authClient *client.AuthClient, eventClient *client.EventClient) *GatewayServer {
	return &GatewayServer{
		authClient:  authClient,
		eventClient: eventClient,
	}
}

func (s *GatewayServer) SignUp(ctx context.Context, req *gateway.SignUpRequest) (*gateway.SignUpResponse, error) {
	resp, err := s.authClient.Register(ctx, &auth.RegisterRequest{
		User: &auth.User{
			Name:      req.GetName(),
			Email:     req.GetEmail(),
			Password:  req.GetPassword(),
			Interests: req.GetInterests(),
		},
	})

	if err != nil {
		return &gateway.SignUpResponse{
			Message: resp.GetMessage(),
		}, err
	}

	return &gateway.SignUpResponse{
		Message: resp.GetMessage(),
	}, err
}

func (s *GatewayServer) SignIn(ctx context.Context, req *gateway.SignInRequest) (*gateway.SignInResponse, error) {
	resp, err := s.authClient.Authenticate(ctx, &auth.AuthenticateRequest{
		Name:     req.GetName(),
		Password: req.GetPassword(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.SignInResponse{
		AccessToken:     resp.GetTokens().GetAccessToken(),
		RefreshToken:    resp.GetTokens().GetRefreshToken(),
		AccessTokenTtl:  resp.GetTokens().GetAccessTokenTtl(),
		RefreshTokenTtl: resp.GetTokens().GetRefreshTokenTtl(),
	}, err
}

func (s *GatewayServer) SignOut(ctx context.Context, req *gateway.SignOutRequest) (*gateway.SignOutResponse, error) {
	return nil, nil
}

func (s *GatewayServer) CreateEvent(ctx context.Context, req *gateway.CreateEventRequest) (*gateway.CreateEventResponse, error) {
	// role creator

	return nil, nil
}

func (s *GatewayServer) ReadEvent(ctx context.Context, req *gateway.ReadEventRequest) (*gateway.ReadEventResponse, error) {
	// role creator

	return nil, nil
}

func (s *GatewayServer) UpdateEvent(ctx context.Context, req *gateway.UpdateEventRequest) (*gateway.UpdateEventResponse, error) {
	// role creator

	return nil, nil
}

func (s *GatewayServer) DeleteEvent(ctx context.Context, req *gateway.DeleteEventRequest) (*gateway.DeleteEventResponse, error) {
	// role creator

	return nil, nil
}

func (s *GatewayServer) ListEvents(ctx context.Context, req *gateway.ListEventsRequest) (*gateway.ListEventsResponse, error) {
	// role creator

	return nil, nil
}

func (s *GatewayServer) ListEventsByCreator(ctx context.Context, req *gateway.ListEventsByCreatorRequest) (*gateway.ListEventsByCreatorResponse, error) {
	// role creator

	return nil, nil
}

func (s *GatewayServer) ListEventsByInterests(ctx context.Context, req *gateway.ListEventsByInterestsRequest) (*gateway.ListEventsByInterestsResponse, error) {
	resp, err := s.eventClient.ListEventsByInterests(ctx, &event.ListEventsByInterestsRequest{
		RequestId: req.GetRequestId(),
		UserId:    req.GetUserId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListEventsByInterestsResponse{
		Events: convertEvents(resp.GetEvents()),
	}, err
}

func (s *GatewayServer) ListRegistratedEvents(ctx context.Context, req *gateway.ListRegistratedEventsRequest) (*gateway.ListRegistratedEventsResponse, error) {
	resp, err := s.eventClient.ListRegistratedEvents(ctx, &event.ListRegistratedEventsRequest{
		RequestId: req.GetRequestId(),
		UserId:    req.GetUserId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListRegistratedEventsResponse{
		Events: convertEvents(resp.GetEvents()),
	}, err
}

func (s *GatewayServer) ListUsersToChat(ctx context.Context, req *gateway.ListUsersToChatRequest) (*gateway.ListUsersToChatResponse, error) {
	resp, err := s.eventClient.ListUsersToChat(ctx, &event.ListUsersToChatRequest{
		RequestId: req.GetRequestId(),
		UserId:    req.GetUserId(),
		EventId:   req.GetEventId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListUsersToChatResponse{
		Participants: convertUsers(resp.GetParticipants()),
	}, err
}

func (s *GatewayServer) RegisterUser(ctx context.Context, req *gateway.RegisterUserRequest) (*gateway.RegisterUserResponse, error) {
	_, err := s.eventClient.RegisterUser(ctx, &event.RegisterUserRequest{
		RequestId: req.GetRequestId(),
		UserId:    req.GetUserId(),
		EventId:   req.GetEventId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.RegisterUserResponse{}, err
}

func (s *GatewayServer) SetChatStatus(ctx context.Context, req *gateway.SetChatStatusRequest) (*gateway.SetChatStatusResponse, error) {
	_, err := s.eventClient.SetChatStatus(ctx, &event.SetChatStatusRequest{
		RequestId:     req.GetRequestId(),
		ParticipantId: req.GetUserId(),
		EventId:       req.GetEventId(),
		IsReady:       req.GetIsReady(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.SetChatStatusResponse{}, err
}
