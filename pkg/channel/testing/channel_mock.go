// Code generated by mockery v1.0.0. DO NOT EDIT.

package testing

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	"github.com/liuxd6825/dapr/pkg/config"
	v1 "github.com/liuxd6825/dapr/pkg/messaging/v1"
)

// MockAppChannel is an autogenerated mock type for the AppChannel type
type MockAppChannel struct {
	mock.Mock
}

// GetAppConfig provides a mock function with given fields:
func (_m *MockAppChannel) GetAppConfig() (*config.ApplicationConfig, error) {
	ret := _m.Called()

	var r0 *config.ApplicationConfig
	if rf, ok := ret.Get(0).(func() *config.ApplicationConfig); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.ApplicationConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBaseAddress provides a mock function with given fields:
func (_m *MockAppChannel) GetBaseAddress() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// InvokeMethod provides a mock function with given fields: ctx, req
func (_m *MockAppChannel) InvokeMethod(ctx context.Context, req *v1.InvokeMethodRequest) (*v1.InvokeMethodResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *v1.InvokeMethodResponse
	if rf, ok := ret.Get(0).(func(context.Context, *v1.InvokeMethodRequest) *v1.InvokeMethodResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.InvokeMethodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *v1.InvokeMethodRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
