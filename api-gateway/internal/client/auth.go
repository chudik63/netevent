package client

import (
	"context"

	"github.com/chudik63/netevent/api-gateway/internal/config"
	auth "github.com/chudik63/netevent/auth/pkg/proto"
	"github.com/chudik63/netevent/event_service/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	auth.AuthServiceClient
	conn *grpc.ClientConn
}

func NewAuthClient(ctx context.Context, cfg *config.Config) *AuthClient {
	logs := logger.GetLoggerFromCtx(ctx)

	conn, err := grpc.NewClient(cfg.AuthServiceHost+":"+cfg.AuthServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logs.Fatal(ctx, "failed to create auth client", zap.String("err", err.Error()))
	}

	client := auth.NewAuthServiceClient(conn)

	return &AuthClient{client, conn}
}

func (c *AuthClient) Stop() error {
	return c.conn.Close()
}
