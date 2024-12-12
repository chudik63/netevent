package grpc

import (
	"context"

	"gitlab.crja72.ru/gospec/go9/netevent/auth/pkg/proto"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	target string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{cfg.AuthServerHost + ":" + cfg.AuthServerPort}
}

func (c *Client) GetUserInterests(ctx context.Context, userID int64) ([]string, error) {
	conn, err := c.getConnector()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	authClient := proto.NewAuthServiceClient(conn)

	response, err := authClient.GetInterests(ctx, &proto.GetInterestsRequest{
		UserId: userID,
	})

	return response.GetInterests(), err
}

func (c *Client) getConnector() (*grpc.ClientConn, error) {
	return grpc.NewClient(c.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
