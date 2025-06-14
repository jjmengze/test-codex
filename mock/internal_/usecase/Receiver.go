// Code generated by mockery v2.53.4. DO NOT EDIT.

package mock

import (
	context "context"
	entity "log-receiver/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Receiver is an autogenerated mock type for the Receiver type
type Receiver struct {
	mock.Mock
}

// PutData provides a mock function with given fields: ctx, input
func (_m *Receiver) PutData(ctx context.Context, input entity.PutDataInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for PutData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.PutDataInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewReceiver creates a new instance of Receiver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReceiver(t interface {
	mock.TestingT
	Cleanup(func())
}) *Receiver {
	mock := &Receiver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
