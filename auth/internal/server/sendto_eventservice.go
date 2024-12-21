package server

import (
	"context"
	"time"

	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/db/postgres/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "event/pkg/api/proto/event"
)

var eventPort = ":5300"

func sendToEvent(data *models.Participant) error {
	conn, err := grpc.NewClient("events:"+eventPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.AddParticipant(ctx, &pb.AddParticipantRequest{
		User: &pb.Participant{
			UserId:    data.UserId,
			Name:      data.Name,
			Interests: data.Interests,
		},
		Email: data.Email,
	})
	if err != nil {
		return err
	}
	return nil
}
