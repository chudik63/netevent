package client

import (
	"context"

	"github.com/chudik63/netevent/api-gateway/internal/config"
	event "github.com/chudik63/netevent/event_service/pkg/api/proto/event"
	"github.com/chudik63/netevent/event_service/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventClient struct {
	event.EventServiceClient
	conn *grpc.ClientConn
}

func NewEventClient(ctx context.Context, cfg *config.Config) *EventClient {
	logs := logger.GetLoggerFromCtx(ctx)

	conn, err := grpc.NewClient(cfg.EventServiceHost+":"+cfg.EventServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logs.Fatal(ctx, "failed to create event client", zap.String("err", err.Error()))
	}

	client := event.NewEventServiceClient(conn)

	return &EventClient{client, conn}
}

func (c *EventClient) Stop() error {
	return c.conn.Close()
}
