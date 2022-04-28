// Code generated by mockery v1.0.0. DO NOT EDIT.

package testing

import (
	"context"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/dapr/dapr/pkg/messaging/v1"
)

// MockDirectMessaging is an autogenerated mock type for the MockDirectMessaging type
type MockDirectMessaging struct {
	mock.Mock
}

// Invoke provides a mock function with given fields: ctx, targetAppID, req
func (_m *MockDirectMessaging) Invoke(ctx context.Context, targetAppID string, req *v1.InvokeMethodRequest) (*v1.InvokeMethodResponse, error) {
	ret := _m.Called(ctx, targetAppID, req)

	var r0 *v1.InvokeMethodResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, *v1.InvokeMethodRequest) *v1.InvokeMethodResponse); ok {
		r0 = rf(ctx, targetAppID, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.InvokeMethodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *v1.InvokeMethodRequest) error); ok {
		r1 = rf(ctx, targetAppID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockDirectMessaging) Close() error {
	return nil
}

type FailingDirectMessaging struct {
	Failure Failure
}

func (f *FailingDirectMessaging) Invoke(ctx context.Context, targetAppID string, req *v1.InvokeMethodRequest) (*v1.InvokeMethodResponse, error) {
	err := f.Failure.PerformFailure(string(req.Message().Data.Value))
	if err != nil {
		return &v1.InvokeMethodResponse{}, err
	}
	return v1.NewInvokeMethodResponse(200, "OK", nil), nil
}