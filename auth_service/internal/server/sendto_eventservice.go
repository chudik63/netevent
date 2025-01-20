package server

import (
	"context"
	"time"

	"github.com/chudik63/netevent/auth_service/internal/db/postgres/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chudik63/netevent/event_service/pkg/api/proto/event"
)

var eventPort = ":5300"

func sendToEvent(data *models.Participant) error {
	conn, err := grpc.NewClient("events"+eventPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c := event.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.AddParticipant(ctx, &event.AddParticipantRequest{
		User: &event.Participant{
			UserId:    data.UserId,
			Name:      data.Name,
			Email:     data.Email,
			Interests: data.Interests,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
