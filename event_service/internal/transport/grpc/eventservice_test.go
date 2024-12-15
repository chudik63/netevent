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
					Time:        "2022-01-01 13:31:00",
					Place:       "test1",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().CreateEvent(gomock.Any(), event).Return(int64(1), nil)
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
					Time:        "2022-01-01 13:31:00",
					Place:       "test2",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().CreateEvent(gomock.Any(), event).Return(int64(2), nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse: &event.CreateEventResponse{
				EventId: 2,
			},
		},
		{
			name: "Internal Error test",
			inputRequest: &event.CreateEventRequest{
				RequestId: "3",
				Event: &event.Event{
					CreatorId:   1,
					Description: "test1",
					Time:        "2022-01-01 13:31:00",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().CreateEvent(gomock.Any(), event).Return(int64(0), errors.New("internal error"))
			},
			expectedStatusCode: codes.Internal,
			expectedResponse:   nil,
		},
		{
			name: "Invalid Argument test",
			inputRequest: &event.CreateEventRequest{
				RequestId: "4",
				Event: &event.Event{
					CreatorId:   1,
					Title:       "test2",
					Description: "test2",
					Time:        "wrong time",
					Place:       "test2",
					Interests:   []string{"test_interest1", "test_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {

			},
			expectedStatusCode: codes.InvalidArgument,
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
			testCase.mockBehavior(service, &models.Event{
				EventID:     testCase.inputRequest.GetEvent().GetEventId(),
				CreatorID:   testCase.inputRequest.GetEvent().GetCreatorId(),
				Title:       testCase.inputRequest.GetEvent().GetTitle(),
				Description: testCase.inputRequest.GetEvent().GetDescription(),
				Time:        testCase.inputRequest.GetEvent().GetTime(),
				Place:       testCase.inputRequest.GetEvent().GetPlace(),
				Topics:      testCase.inputRequest.GetEvent().GetInterests(),
			})

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.CreateEvent(ctx, testCase.inputRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}

func TestReadEvent(t *testing.T) {
	type mockBehavior func(s *mock_grpc.MockService, eventID int64)

	testTable := []struct {
		name               string
		inputReadRequest   *event.ReadEventRequest
		mockBehavior       mockBehavior
		expectedStatusCode codes.Code
		expectedResponse   *event.ReadEventResponse
	}{
		{
			name: "OK test",
			inputReadRequest: &event.ReadEventRequest{
				RequestId: "2",
				EventId:   1,
			},
			mockBehavior: func(s *mock_grpc.MockService, eventID int64) {
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{
					EventID:     eventID,
					CreatorID:   1,
					Title:       "Test Event",
					Description: "Test Description",
					Time:        "2024-12-25T10:00:00Z",
					Place:       "Test Place",
					Topics:      []string{"test1", "test2"},
				}, nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse: &event.ReadEventResponse{
				Event: &event.Event{
					EventId:     1,
					CreatorId:   1,
					Title:       "Test Event",
					Description: "Test Description",
					Time:        "2024-12-25T10:00:00Z",
					Place:       "Test Place",
					Interests:   []string{"test1", "test2"},
				},
			},
		},
		{
			name: "Not Found test",
			inputReadRequest: &event.ReadEventRequest{
				RequestId: "3",
				EventId:   99,
			},
			mockBehavior: func(s *mock_grpc.MockService, eventID int64) {
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(nil, models.ErrWrongEventId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Internal Error test",
			inputReadRequest: &event.ReadEventRequest{
				RequestId: "4",
				EventId:   42,
			},
			mockBehavior: func(s *mock_grpc.MockService, eventID int64) {
				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(nil, errors.New("internal error"))
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
			testCase.mockBehavior(service, testCase.inputReadRequest.GetEventId())

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.ReadEvent(ctx, testCase.inputReadRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	type mockBehavior func(s *mock_grpc.MockService, event *models.Event)

	testTable := []struct {
		name               string
		inputUpdateRequest *event.UpdateEventRequest
		mockBehavior       mockBehavior
		expectedStatusCode codes.Code
		expectedResponse   *event.UpdateEventResponse
	}{
		{
			name: "OK test",
			inputUpdateRequest: &event.UpdateEventRequest{
				RequestId: "2",
				Event: &event.Event{
					EventId:     1,
					CreatorId:   1,
					Title:       "Updated Event",
					Description: "Updated Description",
					Time:        "2024-12-25T10:00:00Z",
					Place:       "Updated Place",
					Interests:   []string{"updated_interest1", "updated_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse:   &event.UpdateEventResponse{},
		},
		{
			name: "Invalid Argument test",
			inputUpdateRequest: &event.UpdateEventRequest{
				RequestId: "3",
				Event: &event.Event{
					EventId:     1,
					CreatorId:   1,
					Title:       "Updated Event",
					Description: "Updated Description",
					Time:        "invalid time",
					Place:       "Updated Place",
					Interests:   []string{"updated_interest1", "updated_interest2"},
				},
			},
			expectedStatusCode: codes.InvalidArgument,
			expectedResponse:   nil,
		},
		{
			name: "Not Found test",
			inputUpdateRequest: &event.UpdateEventRequest{
				RequestId: "4",
				Event: &event.Event{
					EventId:     99, // Non-existing event ID
					CreatorId:   1,
					Title:       "Updated Event",
					Description: "Updated Description",
					Time:        "2024-12-25T10:00:00Z",
					Place:       "Updated Place",
					Interests:   []string{"updated_interest1", "updated_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(models.ErrWrongEventId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Internal Error test",
			inputUpdateRequest: &event.UpdateEventRequest{
				RequestId: "5",
				Event: &event.Event{
					EventId:     42,
					CreatorId:   1,
					Title:       "Updated Event",
					Description: "Updated Description",
					Time:        "2024-12-25T10:00:00Z",
					Place:       "Updated Place",
					Interests:   []string{"updated_interest1", "updated_interest2"},
				},
			},
			mockBehavior: func(s *mock_grpc.MockService, event *models.Event) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(errors.New("internal error"))
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
			testCase.mockBehavior(service, &models.Event{
				EventID:     testCase.inputUpdateRequest.GetEvent().GetEventId(),
				CreatorID:   testCase.inputUpdateRequest.GetEvent().GetCreatorId(),
				Title:       testCase.inputUpdateRequest.GetEvent().GetTitle(),
				Description: testCase.inputUpdateRequest.GetEvent().GetDescription(),
				Time:        testCase.inputUpdateRequest.GetEvent().GetTime(),
				Place:       testCase.inputUpdateRequest.GetEvent().GetPlace(),
				Topics:      testCase.inputUpdateRequest.GetEvent().GetInterests(),
			})

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.UpdateEvent(ctx, testCase.inputUpdateRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	type mockBehavior func(s *mock_grpc.MockService, eventID int64)

	testTable := []struct {
		name               string
		inputDeleteRequest *event.DeleteEventRequest
		mockBehavior       mockBehavior
		expectedStatusCode codes.Code
		expectedResponse   *event.DeleteEventResponse
	}{
		{
			name: "OK test",
			inputDeleteRequest: &event.DeleteEventRequest{
				RequestId: "2",
				EventId:   1,
			},
			mockBehavior: func(s *mock_grpc.MockService, eventID int64) {
				s.EXPECT().DeleteEvent(gomock.Any(), eventID).Return(nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse:   &event.DeleteEventResponse{},
		},
		{
			name: "Not Found test",
			inputDeleteRequest: &event.DeleteEventRequest{
				RequestId: "3",
				EventId:   99,
			},
			mockBehavior: func(s *mock_grpc.MockService, eventID int64) {
				s.EXPECT().DeleteEvent(gomock.Any(), eventID).Return(models.ErrWrongEventId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Internal Error test",
			inputDeleteRequest: &event.DeleteEventRequest{
				RequestId: "4",
				EventId:   42,
			},
			mockBehavior: func(s *mock_grpc.MockService, eventID int64) {
				s.EXPECT().DeleteEvent(gomock.Any(), eventID).Return(errors.New("internal error"))
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
			testCase.mockBehavior(service, testCase.inputDeleteRequest.GetEventId())

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.DeleteEvent(ctx, testCase.inputDeleteRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}

func TestRegisterUser(t *testing.T) {
	type mockBehavior func(s *mock_grpc.MockService, userID, eventID int64)

	testTable := []struct {
		name                 string
		inputRegisterRequest *event.RegisterUserRequest
		mockBehavior         mockBehavior
		expectedStatusCode   codes.Code
		expectedResponse     *event.RegisterUserResponse
	}{
		{
			name: "OK test",
			inputRegisterRequest: &event.RegisterUserRequest{
				RequestId: "1",
				UserId:    1,
				EventId:   1,
			},
			mockBehavior: func(s *mock_grpc.MockService, userID, eventID int64) {
				s.EXPECT().RegisterUser(gomock.Any(), userID, eventID).Return(nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse:   &event.RegisterUserResponse{},
		},
		{
			name: "Not Found test",
			inputRegisterRequest: &event.RegisterUserRequest{
				RequestId: "2",
				UserId:    1,
				EventId:   99,
			},
			mockBehavior: func(s *mock_grpc.MockService, userID, eventID int64) {
				s.EXPECT().RegisterUser(gomock.Any(), userID, eventID).Return(models.ErrWrongEventId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Not Found test",
			inputRegisterRequest: &event.RegisterUserRequest{
				RequestId: "3",
				UserId:    99,
				EventId:   1,
			},
			mockBehavior: func(s *mock_grpc.MockService, userID, eventID int64) {
				s.EXPECT().RegisterUser(gomock.Any(), userID, eventID).Return(models.ErrWrongUserId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Internal Error test",
			inputRegisterRequest: &event.RegisterUserRequest{
				RequestId: "4",
				UserId:    1,
				EventId:   1,
			},
			mockBehavior: func(s *mock_grpc.MockService, userID, eventID int64) {
				s.EXPECT().RegisterUser(gomock.Any(), userID, eventID).Return(errors.New("internal error"))
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
			testCase.mockBehavior(service, testCase.inputRegisterRequest.GetUserId(), testCase.inputRegisterRequest.GetEventId())

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.RegisterUser(ctx, testCase.inputRegisterRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}

func TestSetChatStatus(t *testing.T) {
	type mockBehavior func(s *mock_grpc.MockService, participantID, eventID int64, isReady bool)

	testTable := []struct {
		name               string
		inputRequest       *event.SetChatStatusRequest
		mockBehavior       mockBehavior
		expectedStatusCode codes.Code
		expectedResponse   *event.SetChatStatusResponse
	}{
		{
			name: "OK test",
			inputRequest: &event.SetChatStatusRequest{
				RequestId:     "1",
				ParticipantId: 1,
				EventId:       10,
				IsReady:       true,
			},
			mockBehavior: func(s *mock_grpc.MockService, participantID, eventID int64, isReady bool) {
				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(nil)
			},
			expectedStatusCode: codes.OK,
			expectedResponse:   &event.SetChatStatusResponse{},
		},
		{
			name: "Not Found test",
			inputRequest: &event.SetChatStatusRequest{
				RequestId:     "2",
				ParticipantId: 1,
				EventId:       99,
				IsReady:       true,
			},
			mockBehavior: func(s *mock_grpc.MockService, participantID, eventID int64, isReady bool) {
				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(models.ErrWrongEventId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Not Found test",
			inputRequest: &event.SetChatStatusRequest{
				RequestId:     "3",
				ParticipantId: 99,
				EventId:       10,
				IsReady:       false,
			},
			mockBehavior: func(s *mock_grpc.MockService, participantID, eventID int64, isReady bool) {
				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(models.ErrWrongUserId)
			},
			expectedStatusCode: codes.NotFound,
			expectedResponse:   nil,
		},
		{
			name: "Internal Error test",
			inputRequest: &event.SetChatStatusRequest{
				RequestId:     "4",
				ParticipantId: 1,
				EventId:       10,
				IsReady:       true,
			},
			mockBehavior: func(s *mock_grpc.MockService, participantID, eventID int64, isReady bool) {
				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(errors.New("internal error"))
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
			testCase.mockBehavior(service, testCase.inputRequest.GetParticipantId(), testCase.inputRequest.GetEventId(), testCase.inputRequest.GetIsReady())

			eventservice := NewEventService(ctx, service)

			resp, err := eventservice.SetChatStatus(ctx, testCase.inputRequest)

			assert.Equal(t, testCase.expectedResponse, resp)
			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
		})
	}
}
