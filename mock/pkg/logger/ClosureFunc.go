// Code generated by mockery v2.53.4. DO NOT EDIT.

package mock

import (
	logger "log-receiver/pkg/logger"

	mock "github.com/stretchr/testify/mock"
)

// ClosureFunc is an autogenerated mock type for the ClosureFunc type
type ClosureFunc struct {
	mock.Mock
}

// Execute provides a mock function with no fields
func (_m *ClosureFunc) Execute() []logger.InfLogKeyValue {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []logger.InfLogKeyValue
	if rf, ok := ret.Get(0).(func() []logger.InfLogKeyValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logger.InfLogKeyValue)
		}
	}

	return r0
}

// NewClosureFunc creates a new instance of ClosureFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClosureFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClosureFunc {
	mock := &ClosureFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
