// Code generated by MockGen. DO NOT EDIT.
// Source: event.go

// Package mock_service is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	redis "github.com/redis/go-redis/v9"
	models "github.com/chudik63/netevent/events_service/internal/models"
	producer "github.com/chudik63/netevent/events_service/internal/producer"
	repository "github.com/chudik63/netevent/events_service/internal/repository"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateEvent mocks base method.
func (m *MockRepository) CreateEvent(ctx context.Context, event *models.Event) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", ctx, event)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockRepositoryMockRecorder) CreateEvent(ctx, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockRepository)(nil).CreateEvent), ctx, event)
}

// CreateParticipant mocks base method.
func (m *MockRepository) CreateParticipant(ctx context.Context, participant *models.Participant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateParticipant", ctx, participant)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateParticipant indicates an expected call of CreateParticipant.
func (mr *MockRepositoryMockRecorder) CreateParticipant(ctx, participant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateParticipant", reflect.TypeOf((*MockRepository)(nil).CreateParticipant), ctx, participant)
}

// CreateRegistration mocks base method.
func (m *MockRepository) CreateRegistration(ctx context.Context, userID, eventID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRegistration", ctx, userID, eventID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRegistration indicates an expected call of CreateRegistration.
func (mr *MockRepositoryMockRecorder) CreateRegistration(ctx, userID, eventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRegistration", reflect.TypeOf((*MockRepository)(nil).CreateRegistration), ctx, userID, eventID)
}

// DeleteEvent mocks base method.
func (m *MockRepository) DeleteEvent(ctx context.Context, eventID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, eventID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockRepositoryMockRecorder) DeleteEvent(ctx, eventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockRepository)(nil).DeleteEvent), ctx, eventID)
}

// ListEvents mocks base method.
func (m *MockRepository) ListEvents(ctx context.Context, equations repository.Creds) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEvents", ctx, equations)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEvents indicates an expected call of ListEvents.
func (mr *MockRepositoryMockRecorder) ListEvents(ctx, equations interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEvents", reflect.TypeOf((*MockRepository)(nil).ListEvents), ctx, equations)
}

// ListEventsByInterests mocks base method.
func (m *MockRepository) ListEventsByInterests(ctx context.Context, userID int64, creds repository.Creds) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsByInterests", ctx, userID, creds)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsByInterests indicates an expected call of ListEventsByInterests.
func (mr *MockRepositoryMockRecorder) ListEventsByInterests(ctx, userID, creds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsByInterests", reflect.TypeOf((*MockRepository)(nil).ListEventsByInterests), ctx, userID, creds)
}

// ListRegistratedEvents mocks base method.
func (m *MockRepository) ListRegistratedEvents(ctx context.Context, userID int64) ([]*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRegistratedEvents", ctx, userID)
	ret0, _ := ret[0].([]*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRegistratedEvents indicates an expected call of ListRegistratedEvents.
func (mr *MockRepositoryMockRecorder) ListRegistratedEvents(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRegistratedEvents", reflect.TypeOf((*MockRepository)(nil).ListRegistratedEvents), ctx, userID)
}

// ListUsersToChat mocks base method.
func (m *MockRepository) ListUsersToChat(ctx context.Context, eventID, userID int64) ([]*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsersToChat", ctx, eventID, userID)
	ret0, _ := ret[0].([]*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsersToChat indicates an expected call of ListUsersToChat.
func (mr *MockRepositoryMockRecorder) ListUsersToChat(ctx, eventID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsersToChat", reflect.TypeOf((*MockRepository)(nil).ListUsersToChat), ctx, eventID, userID)
}

// ReadEvent mocks base method.
func (m *MockRepository) ReadEvent(ctx context.Context, eventID int64) (*models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadEvent", ctx, eventID)
	ret0, _ := ret[0].(*models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadEvent indicates an expected call of ReadEvent.
func (mr *MockRepositoryMockRecorder) ReadEvent(ctx, eventID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEvent", reflect.TypeOf((*MockRepository)(nil).ReadEvent), ctx, eventID)
}

// ReadParticipant mocks base method.
func (m *MockRepository) ReadParticipant(ctx context.Context, userID int64) (*models.Participant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadParticipant", ctx, userID)
	ret0, _ := ret[0].(*models.Participant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadParticipant indicates an expected call of ReadParticipant.
func (mr *MockRepositoryMockRecorder) ReadParticipant(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadParticipant", reflect.TypeOf((*MockRepository)(nil).ReadParticipant), ctx, userID)
}

// SetChatStatus mocks base method.
func (m *MockRepository) SetChatStatus(ctx context.Context, userID, eventID int64, isReady bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetChatStatus", ctx, userID, eventID, isReady)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetChatStatus indicates an expected call of SetChatStatus.
func (mr *MockRepositoryMockRecorder) SetChatStatus(ctx, userID, eventID, isReady interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetChatStatus", reflect.TypeOf((*MockRepository)(nil).SetChatStatus), ctx, userID, eventID, isReady)
}

// UpdateEvent mocks base method.
func (m *MockRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockRepositoryMockRecorder) UpdateEvent(ctx, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockRepository)(nil).UpdateEvent), ctx, event)
}

// UpdateParticipant mocks base method.
func (m *MockRepository) UpdateParticipant(ctx context.Context, participant *models.Participant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateParticipant", ctx, participant)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateParticipant indicates an expected call of UpdateParticipant.
func (mr *MockRepositoryMockRecorder) UpdateParticipant(ctx, participant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateParticipant", reflect.TypeOf((*MockRepository)(nil).UpdateParticipant), ctx, participant)
}

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCache) Get(ctx context.Context, key string) *redis.StringCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*redis.StringCmd)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expiration)
	ret0, _ := ret[0].(*redis.StatusCmd)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheMockRecorder) Set(ctx, key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), ctx, key, value, expiration)
}

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Produce mocks base method.
func (m *MockProducer) Produce(ctx context.Context, message producer.Message, topic string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Produce", ctx, message, topic)
}

// Produce indicates an expected call of Produce.
func (mr *MockProducerMockRecorder) Produce(ctx, message, topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockProducer)(nil).Produce), ctx, message, topic)
}
