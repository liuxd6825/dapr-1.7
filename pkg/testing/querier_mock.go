// Code generated by mockery v2.9.4. DO NOT EDIT.

package testing

import (
	state "github.com/liuxd6825/components-contrib/state"

	mock "github.com/stretchr/testify/mock"
)

// MockQuerier is an autogenerated mock type for the Querier type
type MockQuerier struct {
	mock.Mock
}

// Query provides a mock function with given fields: req
func (_m *MockQuerier) Query(req *state.QueryRequest) (*state.QueryResponse, error) {
	ret := _m.Called(req)

	var r0 *state.QueryResponse
	if rf, ok := ret.Get(0).(func(*state.QueryRequest) *state.QueryResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.QueryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*state.QueryRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (f *FailingStatestore) Query(req *state.QueryRequest) (*state.QueryResponse, error) {
	key := req.Metadata["key"]
	err := f.Failure.PerformFailure(key)

	if err != nil {
		return nil, err
	}
	return &state.QueryResponse{}, nil
}
