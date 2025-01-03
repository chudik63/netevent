package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/chudik63/netevent/api-gateway/internal/client"
	"github.com/chudik63/netevent/auth/pkg/proto"

	"github.com/google/uuid"
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
			"/gateway.GatewayService/ListEventsByCreator": true,
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata in request")
		}

		authHeader := md["Authorization"]
		if len(authHeader) == 0 {
			return nil, errors.New("permission denied")
		}

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

		requestID := ""
		if values := md["X-Request-ID"]; len(values) > 0 {
			requestID = values[0]
		} else {
			newUUID, err := uuid.NewUUID()
			if err == nil {
				requestID = newUUID.String()
			}
		}

		ctx = context.WithValue(ctx, "request_id", requestID)

		return handler(ctx, req)
	}
}
