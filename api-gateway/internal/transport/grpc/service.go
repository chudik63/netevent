package grpc

import (
	"context"

	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/internal/client"
	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/pkg/api/gateway"
)

type GatewayServer struct {
	gateway.UnimplementedGatewayServer
	authClient  *client.AuthClient
	eventClient *client.EventClient
}

func NewGatewayServer(authClient *client.AuthClient, eventClient *client.EventClient) *GatewayServer {
	return &GatewayServer{
		authClient:  authClient,
		eventClient: eventClient,
	}
}

func (s *GatewayServer) CreateEvent(ctx context.Context, req *gateway.CreateEventRequest) (*gateway.CreateEventResponse, error) {
	return nil, nil
}

func (s *GatewayServer) DeleteEvent(ctx context.Context, req *gateway.DeleteEventRequest) (*gateway.DeleteEventResponse, error) {
	return nil, nil
}

func (s *GatewayServer) ListEvents(ctx context.Context, req *gateway.ListEventsRequest) (*gateway.ListEventsResponse, error) {
	return nil, nil
}

func (s *GatewayServer) ListEventsByCreator(ctx context.Context, req *gateway.ListEventsByCreatorRequest) (*gateway.ListEventsByCreatorResponse, error) {
	return nil, nil
}

func (s *GatewayServer) ListEventsByInterests(ctx context.Context, req *gateway.ListEventsByInterestsRequest) (*gateway.ListEventsByInterestsResponse, error) {
	return nil, nil
}

func (s *GatewayServer) ListRegistratedEvents(ctx context.Context, req *gateway.ListRegistratedEventsRequest) (*gateway.ListRegistratedEventsResponse, error) {
	return nil, nil
}

func (s *GatewayServer) ListUsersToChat(ctx context.Context, req *gateway.ListUsersToChatRequest) (*gateway.ListUsersToChatResponse, error) {
	return nil, nil
}

func (s *GatewayServer) ReadEvent(ctx context.Context, req *gateway.ReadEventRequest) (*gateway.ReadEventResponse, error) {
	return nil, nil
}

func (s *GatewayServer) RegisterUser(ctx context.Context, req *gateway.RegisterUserRequest) (*gateway.RegisterUserResponse, error) {
	return nil, nil
}

func (s *GatewayServer) SetChatStatus(ctx context.Context, req *gateway.SetChatStatusRequest) (*gateway.SetChatStatusResponse, error) {
	return nil, nil
}

func (s *GatewayServer) SignIn(ctx context.Context, req *gateway.SignInRequest) (*gateway.SignInResponse, error) {
	return nil, nil
}

func (s *GatewayServer) SignOut(ctx context.Context, req *gateway.SignOutRequest) (*gateway.SignOutResponse, error) {
	return nil, nil
}

func (s *GatewayServer) SignUp(ctx context.Context, req *gateway.SignUpRequest) (*gateway.SignUpResponse, error) {
	return nil, nil
}

func (s *GatewayServer) UpdateEvent(ctx context.Context, req *gateway.UpdateEventRequest) (*gateway.UpdateEventResponse, error) {
	return nil, nil
}
