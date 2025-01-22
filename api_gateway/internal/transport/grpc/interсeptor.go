package grpc

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/chudik63/netevent/api_gateway/internal/client"
	"github.com/chudik63/netevent/auth_service/pkg/proto"
	"github.com/chudik63/netevent/events_service/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient *client.AuthClient
	logger     logger.Logger
}

func NewAuthInterceptor(authClient *client.AuthClient, logger logger.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		authClient: authClient,
		logger:     logger,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/gateway.Gateway/SignIn" || info.FullMethod == "/gateway.Gateway/SignUp" {
			return handler(ctx, req)
		}

		creatorRoutes := map[string]bool{
			"/gateway.Gateway/CreateEvent": true,
			"/gateway.Gateway/ReadEvent":   true,
			"/gateway.Gateway/UpdateEvent": true,
			"/gateway.Gateway/DeleteEvent": true,
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata in request")
		}

		authHeader := md.Get("Authorization")
		if len(authHeader) == 0 {
			return nil, errors.New("permission denied")
		}

		token := strings.TrimSpace(authHeader[0])
		if !strings.HasPrefix(token, "Bearer ") {
			return nil, errors.New("invalid authorization header")
		}

		token = strings.TrimPrefix(token, "Bearer ")

		resp, err := i.authClient.Authorise(ctx, &proto.AuthoriseRequest{
			Token: token,
		})
		if err != nil {
			return nil, err
		}

		role := resp.GetRole()
		if role == "" {
			return nil, errors.New("failed to get role from token")
		}

		id := resp.GetId()
		if id == 0 {
			return nil, errors.New("failed to get id from token")
		}

		md.Append("userID", strconv.FormatInt(id, 10))

		if creatorRoutes[info.FullMethod] && role != "creator" {
			return nil, errors.New("permission denied")
		}

		if len(md.Get("X-Request-ID")) == 0 {
			newUUID, err := uuid.NewUUID()
			if err == nil {
				md.Set("X-Request-ID", newUUID.String())
			}
		}

		newCtx := metadata.NewIncomingContext(ctx, md)

		return handler(newCtx, req)
	}
}
