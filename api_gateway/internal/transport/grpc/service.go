package grpc

import (
	"context"

	"github.com/chudik63/netevent/api_gateway/internal/client"
	"github.com/chudik63/netevent/api_gateway/pkg/api/gateway"
	auth "github.com/chudik63/netevent/auth_service/pkg/proto"
	event "github.com/chudik63/netevent/events_service/pkg/api/proto/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	_, err := s.authClient.Register(ctx, &auth.RegisterRequest{
		User: &auth.User{
			Name:      req.GetName(),
			Email:     req.GetEmail(),
			Password:  req.GetPassword(),
			Role:      req.GetRole(),
			Interests: req.GetInterests(),
		},
	})

	if err != nil {
		return nil, err
	}

	return &gateway.SignUpResponse{}, nil
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
	}, nil
}

func (s *GatewayServer) CreateEvent(ctx context.Context, req *gateway.CreateEventRequest) (*gateway.CreateEventResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id in context")
	}

	resp, err := s.eventClient.CreateEvent(ctx, &event.CreateEventRequest{
		Event: &event.Event{
			CreatorId:   userId,
			EventId:     req.GetEvent().GetEventId(),
			Title:       req.GetEvent().GetTitle(),
			Description: req.GetEvent().GetDescription(),
			Time:        req.GetEvent().GetTime(),
			Place:       req.GetEvent().GetPlace(),
			Interests:   req.GetEvent().GetInterests(),
		},
	})
	if err != nil {
		return nil, err
	}

	return &gateway.CreateEventResponse{
		EventId: resp.GetEventId(),
	}, nil
}

func (s *GatewayServer) ReadEvent(ctx context.Context, req *gateway.ReadEventRequest) (*gateway.ReadEventResponse, error) {
	resp, err := s.eventClient.ReadEvent(ctx, &event.ReadEventRequest{
		EventId: req.GetEventId(),
	})
	if err != nil {
		return nil, err
	}

	return &gateway.ReadEventResponse{
		Event: &gateway.Event{
			EventId:     resp.GetEvent().GetEventId(),
			Title:       resp.GetEvent().GetTitle(),
			Description: resp.GetEvent().GetDescription(),
			Time:        resp.GetEvent().GetTime(),
			Place:       resp.GetEvent().GetPlace(),
			Interests:   resp.GetEvent().GetInterests(),
		},
	}, nil
}

func (s *GatewayServer) UpdateEvent(ctx context.Context, req *gateway.UpdateEventRequest) (*gateway.UpdateEventResponse, error) {
	_, err := s.eventClient.UpdateEvent(ctx, &event.UpdateEventRequest{
		Event: &event.Event{
			EventId:     req.GetEvent().GetEventId(),
			CreatorId:   req.GetEvent().GetCreatorId(),
			Title:       req.GetEvent().GetTitle(),
			Description: req.GetEvent().GetDescription(),
			Time:        req.GetEvent().GetTime(),
			Place:       req.GetEvent().GetPlace(),
			Interests:   req.GetEvent().GetInterests(),
		},
	})
	if err != nil {
		return nil, err
	}

	return &gateway.UpdateEventResponse{}, nil
}

func (s *GatewayServer) DeleteEvent(ctx context.Context, req *gateway.DeleteEventRequest) (*gateway.DeleteEventResponse, error) {
	_, err := s.eventClient.DeleteEvent(ctx, &event.DeleteEventRequest{
		EventId: req.GetEventId(),
	})
	if err != nil {
		return nil, err
	}

	return &gateway.DeleteEventResponse{}, nil
}

func (s *GatewayServer) ListEvents(ctx context.Context, req *gateway.ListEventsRequest) (*gateway.ListEventsResponse, error) {
	resp, err := s.eventClient.ListEvents(ctx, &event.ListEventsRequest{})

	if err != nil {
		return nil, err
	}

	return &gateway.ListEventsResponse{
		Events: convertEvents(resp.GetEvents()),
	}, nil
}

func (s *GatewayServer) ListEventsByCreator(ctx context.Context, req *gateway.ListEventsByCreatorRequest) (*gateway.ListEventsByCreatorResponse, error) {
	resp, err := s.eventClient.ListEventsByCreator(ctx, &event.ListEventsByCreatorRequest{
		CreatorId: req.GetCreatorId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListEventsByCreatorResponse{
		Events: convertEvents(resp.GetEvents()),
	}, nil
}

func (s *GatewayServer) ListEventsByInterests(ctx context.Context, req *gateway.ListEventsByInterestsRequest) (*gateway.ListEventsByInterestsResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id in context")
	}

	resp, err := s.eventClient.ListEventsByInterests(ctx, &event.ListEventsByInterestsRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListEventsByInterestsResponse{
		Events: convertEvents(resp.GetEvents()),
	}, nil
}

func (s *GatewayServer) ListRegistratedEvents(ctx context.Context, req *gateway.ListRegistratedEventsRequest) (*gateway.ListRegistratedEventsResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id in context")
	}

	resp, err := s.eventClient.ListRegistratedEvents(ctx, &event.ListRegistratedEventsRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListRegistratedEventsResponse{
		Events: convertEvents(resp.GetEvents()),
	}, nil
}

func (s *GatewayServer) ListUsersToChat(ctx context.Context, req *gateway.ListUsersToChatRequest) (*gateway.ListUsersToChatResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id in context")
	}

	resp, err := s.eventClient.ListUsersToChat(ctx, &event.ListUsersToChatRequest{
		UserId:  userId,
		EventId: req.GetEventId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.ListUsersToChatResponse{
		Participants: convertUsers(resp.GetParticipants()),
	}, nil
}

func (s *GatewayServer) RegisterUser(ctx context.Context, req *gateway.RegisterUserRequest) (*gateway.RegisterUserResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id in context")
	}

	_, err := s.eventClient.RegisterUser(ctx, &event.RegisterUserRequest{
		UserId:  userId,
		EventId: req.GetEventId(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.RegisterUserResponse{}, nil
}

func (s *GatewayServer) SetChatStatus(ctx context.Context, req *gateway.SetChatStatusRequest) (*gateway.SetChatStatusResponse, error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id in context")
	}

	_, err := s.eventClient.SetChatStatus(ctx, &event.SetChatStatusRequest{
		ParticipantId: userId,
		EventId:       req.GetEventId(),
		IsReady:       req.GetIsReady(),
	})

	if err != nil {
		return nil, err
	}

	return &gateway.SetChatStatusResponse{}, nil
}
