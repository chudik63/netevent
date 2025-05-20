package grpc

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/chudik63/netevent/events_service/internal/models"
// 	"github.com/chudik63/netevent/events_service/internal/repository"
// 	"github.com/chudik63/netevent/events_service/internal/transport/grpc/mock"
// 	"github.com/chudik63/netevent/events_service/pkg/api/proto/event"
// 	"github.com/chudik63/netevent/events_service/pkg/logger"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// )

// func TestCreateEvent(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, event *models.Event)

// 	testTable := []struct {
// 		name               string
// 		inputRequest       *event.CreateEventRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.CreateEventResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputRequest: &event.CreateEventRequest{
// 				Event: &event.Event{
// 					CreatorId:   1,
// 					Title:       "test1",
// 					Description: "test1",
// 					Time:        "2022-01-01 13:31:00",
// 					Place:       "test1",
// 					Interests:   []string{"test_interest1", "test_interest2"},
// 				},
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event) {
// 				s.EXPECT().CreateEvent(gomock.Any(), event).Return(int64(1), nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.CreateEventResponse{
// 				EventId: 1,
// 			},
// 		},
// 		{
// 			name: "OK test",
// 			inputRequest: &event.CreateEventRequest{
// 				Event: &event.Event{
// 					CreatorId:   1,
// 					Title:       "test2",
// 					Description: "test2",
// 					Time:        "2022-01-01 13:31:00",
// 					Place:       "test2",
// 					Interests:   []string{"test_interest1", "test_interest2"},
// 				},
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event) {
// 				s.EXPECT().CreateEvent(gomock.Any(), event).Return(int64(2), nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.CreateEventResponse{
// 				EventId: 2,
// 			},
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputRequest: &event.CreateEventRequest{
// 				Event: &event.Event{
// 					CreatorId:   1,
// 					Description: "test1",
// 					Time:        "2022-01-01 13:31:00",
// 					Interests:   []string{"test_interest1", "test_interest2"},
// 				},
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event) {
// 				s.EXPECT().CreateEvent(gomock.Any(), event).Return(int64(0), errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Invalid Argument test",
// 			inputRequest: &event.CreateEventRequest{
// 				Event: &event.Event{
// 					CreatorId:   1,
// 					Title:       "test2",
// 					Description: "test2",
// 					Time:        "wrong time",
// 					Place:       "test2",
// 					Interests:   []string{"test_interest1", "test_interest2"},
// 				},
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event) {

// 			},
// 			expectedStatusCode: codes.InvalidArgument,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, &models.Event{
// 				EventID:     testCase.inputRequest.GetEvent().GetEventId(),
// 				CreatorID:   testCase.inputRequest.GetEvent().GetCreatorId(),
// 				Title:       testCase.inputRequest.GetEvent().GetTitle(),
// 				Description: testCase.inputRequest.GetEvent().GetDescription(),
// 				Time:        testCase.inputRequest.GetEvent().GetTime(),
// 				Place:       testCase.inputRequest.GetEvent().GetPlace(),
// 				Topics:      testCase.inputRequest.GetEvent().GetInterests(),
// 			})

// 			eventservice := NewEventService(ctx, service)

// 			resp, err := eventservice.CreateEvent(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestReadEvent(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, eventID int64)

// 	testTable := []struct {
// 		name               string
// 		inputReadRequest   *event.ReadEventRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.ReadEventResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputReadRequest: &event.ReadEventRequest{
// 				EventId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64) {
// 				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(&models.Event{
// 					EventID:     eventID,
// 					CreatorID:   1,
// 					Title:       "Test Event",
// 					Description: "Test Description",
// 					Time:        "2024-12-25T10:00:00Z",
// 					Place:       "Test Place",
// 					Topics:      []string{"test1", "test2"},
// 				}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ReadEventResponse{
// 				Event: &event.Event{
// 					EventId:     1,
// 					CreatorId:   1,
// 					Title:       "Test Event",
// 					Description: "Test Description",
// 					Time:        "2024-12-25T10:00:00Z",
// 					Place:       "Test Place",
// 					Interests:   []string{"test1", "test2"},
// 				},
// 			},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputReadRequest: &event.ReadEventRequest{
// 				EventId: 99,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64) {
// 				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(nil, models.ErrWrongEventId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputReadRequest: &event.ReadEventRequest{
// 				EventId: 42,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64) {
// 				s.EXPECT().ReadEvent(gomock.Any(), eventID).Return(nil, errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, testCase.inputReadRequest.GetEventId())

// 			eventservice := NewEventService(ctx, service)

// 			resp, err := eventservice.ReadEvent(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputReadRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestUpdateEvent(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, event *models.Event, userId int64)

// 	testTable := []struct {
// 		name               string
// 		inputUpdateRequest *event.UpdateEventRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.UpdateEventResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputUpdateRequest: &event.UpdateEventRequest{
// 				Event: &event.Event{
// 					EventId:     1,
// 					CreatorId:   1,
// 					Title:       "Updated Event",
// 					Description: "Updated Description",
// 					Time:        "2024-01-01 10:00:00",
// 					Place:       "Updated Place",
// 					Interests:   []string{"updated_interest1", "updated_interest2"},
// 				},
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event, userId int64) {
// 				s.EXPECT().UpdateEvent(gomock.Any(), event, userId).Return(nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse:   &event.UpdateEventResponse{},
// 		},
// 		{
// 			name: "Invalid Argument test",
// 			inputUpdateRequest: &event.UpdateEventRequest{
// 				Event: &event.Event{
// 					EventId:     1,
// 					CreatorId:   1,
// 					Title:       "Updated Event",
// 					Description: "Updated Description",
// 					Time:        "invalid time",
// 					Place:       "Updated Place",
// 					Interests:   []string{"updated_interest1", "updated_interest2"},
// 				},
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event, userId int64) {
// 			},
// 			expectedStatusCode: codes.InvalidArgument,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Not Found test",
// 			inputUpdateRequest: &event.UpdateEventRequest{
// 				Event: &event.Event{
// 					EventId:     99,
// 					CreatorId:   1,
// 					Title:       "Updated Event",
// 					Description: "Updated Description",
// 					Time:        "2024-01-01 10:00:00",
// 					Place:       "Updated Place",
// 					Interests:   []string{"updated_interest1", "updated_interest2"},
// 				},
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event, userId int64) {
// 				s.EXPECT().UpdateEvent(gomock.Any(), event, userId).Return(models.ErrWrongEventId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputUpdateRequest: &event.UpdateEventRequest{
// 				Event: &event.Event{
// 					EventId:     42,
// 					CreatorId:   1,
// 					Title:       "Updated Event",
// 					Description: "Updated Description",
// 					Time:        "2024-01-01 10:00:00",
// 					Place:       "Updated Place",
// 					Interests:   []string{"updated_interest1", "updated_interest2"},
// 				},
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event, userId int64) {
// 				s.EXPECT().UpdateEvent(gomock.Any(), event, userId).Return(errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Access denied error test",
// 			inputUpdateRequest: &event.UpdateEventRequest{
// 				Event: &event.Event{
// 					EventId:     42,
// 					CreatorId:   1,
// 					Title:       "Updated Event",
// 					Description: "Updated Description",
// 					Time:        "2024-01-01 10:00:00",
// 					Place:       "Updated Place",
// 					Interests:   []string{"updated_interest1", "updated_interest2"},
// 				},
// 				UserId: 11,
// 			},
// 			mockBehavior: func(s *mock.MockService, event *models.Event, userId int64) {
// 				s.EXPECT().UpdateEvent(gomock.Any(), event, userId).Return(models.ErrAccessDenied)
// 			},
// 			expectedStatusCode: codes.PermissionDenied,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, &models.Event{
// 				EventID:     testCase.inputUpdateRequest.GetEvent().GetEventId(),
// 				CreatorID:   testCase.inputUpdateRequest.GetEvent().GetCreatorId(),
// 				Title:       testCase.inputUpdateRequest.GetEvent().GetTitle(),
// 				Description: testCase.inputUpdateRequest.GetEvent().GetDescription(),
// 				Time:        testCase.inputUpdateRequest.GetEvent().GetTime(),
// 				Place:       testCase.inputUpdateRequest.GetEvent().GetPlace(),
// 				Topics:      testCase.inputUpdateRequest.GetEvent().GetInterests(),
// 			}, testCase.inputUpdateRequest.GetUserId())

// 			eventservice := NewEventService(ctx, service)

// 			resp, err := eventservice.UpdateEvent(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputUpdateRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestDeleteEvent(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, eventID int64, userId int64)

// 	testTable := []struct {
// 		name               string
// 		inputDeleteRequest *event.DeleteEventRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.DeleteEventResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputDeleteRequest: &event.DeleteEventRequest{
// 				EventId: 1,
// 				UserId:  1,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userId int64) {
// 				s.EXPECT().DeleteEvent(gomock.Any(), eventID, userId).Return(nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse:   &event.DeleteEventResponse{},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputDeleteRequest: &event.DeleteEventRequest{
// 				EventId: 99,
// 				UserId:  1,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userId int64) {
// 				s.EXPECT().DeleteEvent(gomock.Any(), eventID, userId).Return(models.ErrWrongEventId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputDeleteRequest: &event.DeleteEventRequest{
// 				EventId: 42,
// 				UserId:  1,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userId int64) {
// 				s.EXPECT().DeleteEvent(gomock.Any(), eventID, userId).Return(errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Access Denied Error test",
// 			inputDeleteRequest: &event.DeleteEventRequest{
// 				EventId: 42,
// 				UserId:  1211,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userId int64) {
// 				s.EXPECT().DeleteEvent(gomock.Any(), eventID, userId).Return(models.ErrAccessDenied)
// 			},
// 			expectedStatusCode: codes.PermissionDenied,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, testCase.inputDeleteRequest.GetEventId(), testCase.inputDeleteRequest.GetUserId())

// 			eventservice := NewEventService(ctx, service)

// 			resp, err := eventservice.DeleteEvent(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputDeleteRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestListEvents(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, creds repository.Creds)

// 	testTable := []struct {
// 		name               string
// 		inputRequest       *event.ListEventsRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.ListEventsResponse
// 	}{
// 		{
// 			name:         "OK test",
// 			inputRequest: &event.ListEventsRequest{},
// 			mockBehavior: func(s *mock.MockService, creds repository.Creds) {
// 				s.EXPECT().ListEvents(gomock.Any(), creds).Return([]*models.Event{{EventID: 1, CreatorID: 1}, {EventID: 2, CreatorID: 2}}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ListEventsResponse{
// 				Events: []*event.Event{{EventId: 1, CreatorId: 1}, {EventId: 2, CreatorId: 2}},
// 			},
// 		},
// 		{
// 			name:         "With creator id test",
// 			inputRequest: &event.ListEventsRequest{CreatorId: 1},
// 			mockBehavior: func(s *mock.MockService, creds repository.Creds) {
// 				s.EXPECT().ListEvents(gomock.Any(), creds).Return([]*models.Event{{EventID: 1, CreatorID: 1}}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ListEventsResponse{
// 				Events: []*event.Event{{EventId: 1, CreatorId: 1}},
// 			},
// 		},
// 		{
// 			name:         "Error not found test",
// 			inputRequest: &event.ListEventsRequest{CreatorId: 11212},
// 			mockBehavior: func(s *mock.MockService, creds repository.Creds) {
// 				s.EXPECT().ListEvents(gomock.Any(), creds).Return(nil, models.ErrNotFound)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)

// 			creds := repository.Creds{}
// 			if creatorID := testCase.inputRequest.GetCreatorId(); creatorID != 0 {
// 				creds["creator_id"] = creatorID
// 			}
// 			testCase.mockBehavior(service, creds)

// 			eventService := NewEventService(ctx, service)

// 			resp, err := eventService.ListEvents(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestListEventsByInterests(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, userId int64, creds repository.Creds)

// 	testTable := []struct {
// 		name               string
// 		inputRequest       *event.ListEventsByInterestsRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.ListEventsByInterestsResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputRequest: &event.ListEventsByInterestsRequest{
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userId int64, creds repository.Creds) {
// 				s.EXPECT().ListEventsByInterests(gomock.Any(), userId, creds).Return([]*models.Event{{EventID: 1, CreatorID: 1}, {EventID: 2, CreatorID: 2}}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ListEventsByInterestsResponse{
// 				Events: []*event.Event{{EventId: 1, CreatorId: 1}, {EventId: 2, CreatorId: 2}},
// 			},
// 		},
// 		{
// 			name: "User not found test",
// 			inputRequest: &event.ListEventsByInterestsRequest{
// 				UserId: 99,
// 			},
// 			mockBehavior: func(s *mock.MockService, userId int64, creds repository.Creds) {
// 				s.EXPECT().ListEventsByInterests(gomock.Any(), userId, creds).Return(nil, models.ErrWrongUserId)
// 			},
// 			expectedStatusCode: codes.InvalidArgument,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "With creator id filter",
// 			inputRequest: &event.ListEventsByInterestsRequest{
// 				UserId:    1,
// 				CreatorId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userId int64, creds repository.Creds) {
// 				s.EXPECT().ListEventsByInterests(gomock.Any(), userId, creds).Return([]*models.Event{{EventID: 1, CreatorID: 1}}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ListEventsByInterestsResponse{
// 				Events: []*event.Event{{EventId: 1, CreatorId: 1}},
// 			},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRequest: &event.ListEventsByInterestsRequest{
// 				UserId:    1,
// 				CreatorId: 3232313,
// 			},
// 			mockBehavior: func(s *mock.MockService, userId int64, creds repository.Creds) {
// 				s.EXPECT().ListEventsByInterests(gomock.Any(), userId, creds).Return(nil, models.ErrNotFound)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)

// 			creds := repository.Creds{}
// 			if creatorID := testCase.inputRequest.GetCreatorId(); creatorID != 0 {
// 				creds["creator_id"] = creatorID
// 			}
// 			testCase.mockBehavior(service, testCase.inputRequest.GetUserId(), creds)

// 			eventService := NewEventService(ctx, service)

// 			resp, err := eventService.ListEventsByInterests(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestRegisterUser(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, userID, eventID int64)

// 	testTable := []struct {
// 		name                 string
// 		inputRegisterRequest *event.RegisterUserRequest
// 		mockBehavior         mockBehavior
// 		expectedStatusCode   codes.Code
// 		expectedResponse     *event.RegisterUserResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputRegisterRequest: &event.RegisterUserRequest{
// 				UserId:  1,
// 				EventId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID, eventID int64) {
// 				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse:   &event.RegisterUserResponse{},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRegisterRequest: &event.RegisterUserRequest{
// 				UserId:  1,
// 				EventId: 99,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID, eventID int64) {
// 				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(models.ErrWrongEventId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRegisterRequest: &event.RegisterUserRequest{
// 				UserId:  99,
// 				EventId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID, eventID int64) {
// 				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(models.ErrWrongUserId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputRegisterRequest: &event.RegisterUserRequest{
// 				UserId:  1,
// 				EventId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID, eventID int64) {
// 				s.EXPECT().CreateRegistration(gomock.Any(), userID, eventID).Return(errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, testCase.inputRegisterRequest.GetUserId(), testCase.inputRegisterRequest.GetEventId())

// 			eventservice := NewEventService(ctx, service)

// 			resp, err := eventservice.RegisterUser(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRegisterRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestSetChatStatus(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, participantID, eventID int64, isReady bool)

// 	testTable := []struct {
// 		name               string
// 		inputRequest       *event.SetChatStatusRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.SetChatStatusResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputRequest: &event.SetChatStatusRequest{
// 				ParticipantId: 1,
// 				EventId:       10,
// 				IsReady:       true,
// 			},
// 			mockBehavior: func(s *mock.MockService, participantID, eventID int64, isReady bool) {
// 				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse:   &event.SetChatStatusResponse{},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRequest: &event.SetChatStatusRequest{
// 				ParticipantId: 1,
// 				EventId:       99,
// 				IsReady:       true,
// 			},
// 			mockBehavior: func(s *mock.MockService, participantID, eventID int64, isReady bool) {
// 				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(models.ErrWrongEventId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRequest: &event.SetChatStatusRequest{
// 				ParticipantId: 99,
// 				EventId:       10,
// 				IsReady:       false,
// 			},
// 			mockBehavior: func(s *mock.MockService, participantID, eventID int64, isReady bool) {
// 				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(models.ErrWrongUserId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputRequest: &event.SetChatStatusRequest{
// 				ParticipantId: 1,
// 				EventId:       10,
// 				IsReady:       true,
// 			},
// 			mockBehavior: func(s *mock.MockService, participantID, eventID int64, isReady bool) {
// 				s.EXPECT().SetChatStatus(gomock.Any(), participantID, eventID, isReady).Return(errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, testCase.inputRequest.GetParticipantId(), testCase.inputRequest.GetEventId(), testCase.inputRequest.GetIsReady())

// 			eventservice := NewEventService(ctx, service)

// 			resp, err := eventservice.SetChatStatus(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestListRegistratedEvents(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, userID int64)

// 	testTable := []struct {
// 		name               string
// 		inputRequest       *event.ListRegistratedEventsRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.ListRegistratedEventsResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputRequest: &event.ListRegistratedEventsRequest{
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID int64) {
// 				s.EXPECT().ListRegistratedEvents(gomock.Any(), userID).Return([]*models.Event{{EventID: 1}}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ListRegistratedEventsResponse{
// 				Events: []*event.Event{{EventId: 1}},
// 			},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRequest: &event.ListRegistratedEventsRequest{
// 				UserId: 99,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID int64) {
// 				s.EXPECT().ListRegistratedEvents(gomock.Any(), userID).Return(nil, models.ErrWrongUserId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputRequest: &event.ListRegistratedEventsRequest{
// 				UserId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, userID int64) {
// 				s.EXPECT().ListRegistratedEvents(gomock.Any(), userID).Return(nil, errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, testCase.inputRequest.GetUserId())

// 			eventService := NewEventService(ctx, service)

// 			resp, err := eventService.ListRegistratedEvents(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }

// func TestListUsersToChat(t *testing.T) {
// 	type mockBehavior func(s *mock.MockService, eventID int64, userID int64)

// 	testTable := []struct {
// 		name               string
// 		inputRequest       *event.ListUsersToChatRequest
// 		mockBehavior       mockBehavior
// 		expectedStatusCode codes.Code
// 		expectedResponse   *event.ListUsersToChatResponse
// 	}{
// 		{
// 			name: "OK test",
// 			inputRequest: &event.ListUsersToChatRequest{
// 				EventId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userID int64) {
// 				s.EXPECT().ListUsersToChat(gomock.Any(), eventID, userID).Return([]*models.Participant{{UserID: 1, Name: "User1"}}, nil)
// 			},
// 			expectedStatusCode: codes.OK,
// 			expectedResponse: &event.ListUsersToChatResponse{
// 				Participants: []*event.Participant{
// 					{UserId: 1, Name: "User1"},
// 				},
// 			},
// 		},
// 		{
// 			name: "Not Found test",
// 			inputRequest: &event.ListUsersToChatRequest{
// 				EventId: 99,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userID int64) {
// 				s.EXPECT().ListUsersToChat(gomock.Any(), eventID, userID).Return(nil, models.ErrWrongEventId)
// 			},
// 			expectedStatusCode: codes.NotFound,
// 			expectedResponse:   nil,
// 		},
// 		{
// 			name: "Internal Error test",
// 			inputRequest: &event.ListUsersToChatRequest{
// 				EventId: 1,
// 			},
// 			mockBehavior: func(s *mock.MockService, eventID int64, userID int64) {
// 				s.EXPECT().ListUsersToChat(gomock.Any(), eventID, userID).Return(nil, errors.New("internal error"))
// 			},
// 			expectedStatusCode: codes.Internal,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			log, _ := logger.New("test")
// 			ctx := context.WithValue(context.Background(), logger.LoggerKey, log)

// 			service := mock.NewMockService(c)
// 			testCase.mockBehavior(service, testCase.inputRequest.GetEventId(), testCase.inputRequest.GetUserId())

// 			eventService := NewEventService(ctx, service)

// 			resp, err := eventService.ListUsersToChat(metadata.NewIncomingContext(ctx, metadata.MD{"x-request-id": []string{"test"}}), testCase.inputRequest)

// 			assert.Equal(t, testCase.expectedResponse, resp)
// 			assert.Equal(t, testCase.expectedStatusCode, status.Code(err))
// 		})
// 	}
// }
