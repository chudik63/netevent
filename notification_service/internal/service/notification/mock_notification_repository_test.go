// Code generated by mockery v2.45.1. DO NOT EDIT.

package notification_test

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/chudik63/netevent/notification_service/internal/domain"
)

// MockNotificationRepository is an autogenerated mock type for the NotificationRepository type
type MockNotificationRepository struct {
	mock.Mock
}

// AddNotification provides a mock function with given fields: ctx, notify
func (_m *MockNotificationRepository) AddNotification(ctx context.Context, notify domain.Notification) (domain.Notification, error) {
	ret := _m.Called(ctx, notify)

	if len(ret) == 0 {
		panic("no return value specified for AddNotification")
	}

	var r0 domain.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Notification) (domain.Notification, error)); ok {
		return rf(ctx, notify)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Notification) domain.Notification); ok {
		r0 = rf(ctx, notify)
	} else {
		r0 = ret.Get(0).(domain.Notification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Notification) error); ok {
		r1 = rf(ctx, notify)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNotification provides a mock function with given fields: ctx, id
func (_m *MockNotificationRepository) DeleteNotification(ctx context.Context, id int64) (domain.Notification, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNotification")
	}

	var r0 domain.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (domain.Notification, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Notification); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Notification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNearestNotifications provides a mock function with given fields: ctx
func (_m *MockNotificationRepository) GetNearestNotifications(ctx context.Context) ([]domain.Notification, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetNearestNotifications")
	}

	var r0 []domain.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.Notification, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Notification); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Notification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockNotificationRepository creates a new instance of MockNotificationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNotificationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNotificationRepository {
	mock := &MockNotificationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
