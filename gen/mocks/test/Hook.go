// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	logrus "github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
)

// Hook is an autogenerated mock type for the Hook type
type Hook struct {
	mock.Mock
}

// Fire provides a mock function with given fields: _a0
func (_m *Hook) Fire(_a0 *logrus.Entry) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*logrus.Entry) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Levels provides a mock function with given fields:
func (_m *Hook) Levels() []logrus.Level {
	ret := _m.Called()

	var r0 []logrus.Level
	if rf, ok := ret.Get(0).(func() []logrus.Level); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logrus.Level)
		}
	}

	return r0
}

type mockConstructorTestingTNewHook interface {
	mock.TestingT
	Cleanup(func())
}

// NewHook creates a new instance of Hook. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHook(t mockConstructorTestingTNewHook) *Hook {
	mock := &Hook{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
