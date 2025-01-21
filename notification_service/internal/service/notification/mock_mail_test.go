// Code generated by mockery v2.45.1. DO NOT EDIT.

package notification_test

import (
	mock "github.com/stretchr/testify/mock"
	domain "github.com/chudik63/netevent/notification_service/internal/domain"
)

// MockMail is an autogenerated mock type for the Mail type
type MockMail struct {
	mock.Mock
}

// Send provides a mock function with given fields: subject, msg
func (_m *MockMail) Send(subject string, msg domain.Notification) error {
	ret := _m.Called(subject, msg)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, domain.Notification) error); ok {
		r0 = rf(subject, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockMail creates a new instance of MockMail. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMail(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMail {
	mock := &MockMail{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
