// Code generated by MockGen. DO NOT EDIT.
// Source: eventservice.go

// Package mock_grpc is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/chudik63/netevent/events_service/internal/models"
	event "github.com/chudik63/netevent/events_service/pkg/api/proto/event"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddParticipant mocks base method.
func (m *MockService) AddParticipant(ctx context.Context, req *event.AddParticipantRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddParticipant", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddParticipant indicates an expected call of AddParticipant.
func (mr *MockServiceMockRecorder) AddParticipant(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddParticipant", reflect.TypeOf((*MockService)(nil).AddParticipant), ctx, req)
}

// CreateEvent mocks base method.
func (m *MockService) CreateEvent(ctx context.Context, req *event.CreateEventRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", ctx, req)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockServiceMockRecorder) CreateEvent(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockService)(nil).CreateEvent), ctx, req)
}

// CreateRegistration mocks base method.
func (m *MockService) CreateRegistration(ctx context.Context, req *event.RegisterUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRegistration", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRegistration indicates an expected call of CreateRegistration.
func (mr *MockServiceMockRecorder) CreateRegistration(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRegistration", reflect.TypeOf((*MockService)(nil).CreateRegistration), ctx, req)
}

// DeleteEvent mocks base method.
func (m *MockService) DeleteEvent(ctx context.Context, req *event.DeleteEventRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockServiceMockRecorder) DeleteEvent(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockService)(nil).DeleteEvent), ctx, req)
}

// ListEvents mocks base method.
func (m *MockService) ListEvents(ctx context.Context, creatorID int64) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEvents", ctx, creatorID)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEvents indicates an expected call of ListEvents.
func (mr *MockServiceMockRecorder) ListEvents(ctx, creatorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEvents", reflect.TypeOf((*MockService)(nil).ListEvents), ctx, creatorID)
}

// ListEventsByInterests mocks base method.
func (m *MockService) ListEventsByInterests(ctx context.Context, req *event.ListEventsByInterestsRequest) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsByInterests", ctx, req)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsByInterests indicates an expected call of ListEventsByInterests.
func (mr *MockServiceMockRecorder) ListEventsByInterests(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsByInterests", reflect.TypeOf((*MockService)(nil).ListEventsByInterests), ctx, req)
}

// ListRegistratedEvents mocks base method.
func (m *MockService) ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRegistratedEvents", ctx, userID)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRegistratedEvents indicates an expected call of ListRegistratedEvents.
func (mr *MockServiceMockRecorder) ListRegistratedEvents(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRegistratedEvents", reflect.TypeOf((*MockService)(nil).ListRegistratedEvents), ctx, userID)
}

// ListUsersToChat mocks base method.
func (m *MockService) ListUsersToChat(ctx context.Context, req *event.ListUsersToChatRequest) ([]*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsersToChat", ctx, req)
	ret0, _ := ret[0].([]*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsersToChat indicates an expected call of ListUsersToChat.
func (mr *MockServiceMockRecorder) ListUsersToChat(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsersToChat", reflect.TypeOf((*MockService)(nil).ListUsersToChat), ctx, req)
}

// ReadEvent mocks base method.
func (m *MockService) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadEvent", ctx, eventID)
	ret0, _ := ret[0].(*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadEvent indicates an expected call of ReadEvent.
func (mr *MockServiceMockRecorder) ReadEvent(ctx, eventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEvent", reflect.TypeOf((*MockService)(nil).ReadEvent), ctx, eventID)
}

// SetChatStatus mocks base method.
func (m *MockService) SetChatStatus(ctx context.Context, req *event.SetChatStatusRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetChatStatus", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetChatStatus indicates an expected call of SetChatStatus.
func (mr *MockServiceMockRecorder) SetChatStatus(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetChatStatus", reflect.TypeOf((*MockService)(nil).SetChatStatus), ctx, req)
}

// UpdateEvent mocks base method.
func (m *MockService) UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockServiceMockRecorder) UpdateEvent(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockService)(nil).UpdateEvent), ctx, req)
}
