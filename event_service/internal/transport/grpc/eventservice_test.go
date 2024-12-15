package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/logger"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/models"
	mock_grpc "gitlab.crja72.ru/gospec/go9/netevent/event_service/internal/transport/grpc/mocks"
	"gitlab.crja72.ru/gospec/go9/netevent/event_service/pkg/api/proto/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateEvent(t *testing.T) {
	type mockBehavior func(s *mock_grpc.MockService, event *models.Event)

	testTable := []struct {
		name               string
		inputRequest       *event.CreateEventRequest
		inputEvent         *models.Event
		mockBehavior       mockBehavior
		expectedStatusCode codes.Code
		expectedResponse   *event.CreateEventResponse
	}{
		{
			name: "OK test",
			inputRequest: &event.CreateEventRequest{
				RequestId: "1",
				Event: &event.Event{
					CreatorId:   1,
					Title:       "test1",
					Description: "test1",
					Time:        "2022-01-01 13:31",
					Place:       "test1",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			inputEvent: &models.Event{
				CreatorID:   1,
				Title:       "test1",
				Description: "test1",
				Time:        "2022-01-01 13:31",
				Place:       "test1",
				Topics:      []string{"test_interest1", "test_interest2"},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().CreateEvent(context.Background(), event).Return(int64(1), nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse: &event.CreateEventResponse{
				EventId: 1,
			},
		},
		{
			name: "OK test",
			inputRequest: &event.CreateEventRequest{
				RequestId: "2",
				Event: &event.Event{
					CreatorId:   1,
					Title:       "test2",
					Description: "test2",
					Time:        "2022-01-01 13:31",
					Place:       "test2",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			inputEvent: &models.Event{
				CreatorID:   1,
				Title:       "test2",
				Description: "test2",
				Time:        "2022-01-01 13:31",
				Place:       "test2",
				Topics:      []string{"test_interest1", "test_interest2"},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().CreateEvent(context.Background(), event).Return(int64(2), nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse: &event.CreateEventResponse{
				EventId: 2,
			},
		},
		{
			name: "Internal Error test",
			inputRequest: &event.CreateEventRequest{
				RequestId: "1",
				Event: &event.Event{
					CreatorId:   1,
					Description: "test1",
					Time:        "2022-01-01 13:31",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			inputEvent: &models.Event{
				CreatorID:   1,
				Description: "test1",
				Time:        "2022-01-01 13:31",
				Topics:      []string{"test_interest1", "test_interest2"},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().CreateEvent(context.Background(), event).Return(int64(0), errors.New("internal error"))
			},
			expectedStatusCode: codes.Internal,
			expectedResponse:   nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			log, _ := logger.New("test")
			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

			service := mock_grpc.NewMockService(c)
			testCase.mockBehavior(service, testCase.inputEvent)

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.CreateEvent(context.Background(), testCase.inputRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}
