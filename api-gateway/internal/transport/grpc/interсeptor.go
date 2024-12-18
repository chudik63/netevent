package grpc

import (
	"context"
	"errors"
	"strings"

	"gitlab.crja72.ru/gospec/go9/netevent/api-gateway/internal/client"
	"gitlab.crja72.ru/gospec/go9/netevent/auth/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient *client.AuthClient
}

func NewAuthInterceptor(authClient *client.AuthClient) *AuthInterceptor {
	return &AuthInterceptor{authClient: authClient}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/gateway.GatewayService/SignIn" || info.FullMethod == "/gateway.GatewayService/SignUp" || info.FullMethod == "/gateway.GatewayService/SignOut" {
			return handler(ctx, req)
		}

		creatorRoutes := map[string]bool{
			"/gateway.GatewayService/CreateEvent":         true,
			"/gateway.GatewayService/ReadEvent":           true,
			"/gateway.GatewayService/UpdateEvent":         true,
			"/gateway.GatewayService/DeleteEvent":         true,
			"/gateway.GatewayService/ListEvents":          true,
			"/gateway.GatewayService/ListEventsByCreator": true,
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata in request")
		}

		authHeader := md["authorization"]
		if len(authHeader) > 0 {
			token := strings.TrimSpace(authHeader[0])
			if strings.HasPrefix(token, "Bearer ") {
				token = strings.TrimPrefix(token, "Bearer ")

				resp, err := i.authClient.Authorise(ctx, &proto.AuthoriseRequest{
					Token: token,
				})
				if err != nil {
					return resp.GetMessage(), err
				}

				role := resp.GetRole()

				if creatorRoutes[info.FullMethod] && role != "creator" {
					return nil, errors.New("permission denied")
				}
			}
		}

		return handler(ctx, req)
	}
}
