// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	adapters "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	client "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"

	config "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Chain is an autogenerated mock type for the Chain type
type Chain struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Chain) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Config provides a mock function with given fields:
func (_m *Chain) Config() config.Config {
	ret := _m.Called()

	var r0 config.Config
	if rf, ok := ret.Get(0).(func() config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.Config)
		}
	}

	return r0
}

// HealthReport provides a mock function with given fields:
func (_m *Chain) HealthReport() map[string]error {
	ret := _m.Called()

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *Chain) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Chain) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Reader provides a mock function with given fields: nodeName
func (_m *Chain) Reader(nodeName string) (client.Reader, error) {
	ret := _m.Called(nodeName)

	var r0 client.Reader
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (client.Reader, error)); ok {
		return rf(nodeName)
	}
	if rf, ok := ret.Get(0).(func(string) client.Reader); ok {
		r0 = rf(nodeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Reader)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(nodeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ready provides a mock function with given fields:
func (_m *Chain) Ready() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *Chain) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TxManager provides a mock function with given fields:
func (_m *Chain) TxManager() adapters.TxManager {
	ret := _m.Called()

	var r0 adapters.TxManager
	if rf, ok := ret.Get(0).(func() adapters.TxManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(adapters.TxManager)
		}
	}

	return r0
}

type mockConstructorTestingTNewChain interface {
	mock.TestingT
	Cleanup(func())
}

// NewChain creates a new instance of Chain. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewChain(t mockConstructorTestingTNewChain) *Chain {
	mock := &Chain{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
