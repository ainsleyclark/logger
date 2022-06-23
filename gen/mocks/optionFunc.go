// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	logger "github.com/krang-backlink/logger"
	mock "github.com/stretchr/testify/mock"
)

// OptionFunc is an autogenerated mock type for the optionFunc type
type OptionFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: config
func (_m *OptionFunc) Execute(config *logger.Config) {
	_m.Called(config)
}

type mockConstructorTestingTNewOptionFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewOptionFunc creates a new instance of OptionFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOptionFunc(t mockConstructorTestingTNewOptionFunc) *OptionFunc {
	mock := &OptionFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
