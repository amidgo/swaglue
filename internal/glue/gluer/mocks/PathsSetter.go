// Code generated by mockery v2.40.1. DO NOT EDIT.

package gluermocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// PathsSetter is an autogenerated mock type for the PathsSetter type
type PathsSetter struct {
	mock.Mock
}

type PathsSetter_Expecter struct {
	mock *mock.Mock
}

func (_m *PathsSetter) EXPECT() *PathsSetter_Expecter {
	return &PathsSetter_Expecter{mock: &_m.Mock}
}

// SetPaths provides a mock function with given fields: paths
func (_m *PathsSetter) SetPaths(paths map[string]io.Reader) error {
	ret := _m.Called(paths)

	if len(ret) == 0 {
		panic("no return value specified for SetPaths")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]io.Reader) error); ok {
		r0 = rf(paths)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PathsSetter_SetPaths_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetPaths'
type PathsSetter_SetPaths_Call struct {
	*mock.Call
}

// SetPaths is a helper method to define mock.On call
//   - paths map[string]io.Reader
func (_e *PathsSetter_Expecter) SetPaths(paths interface{}) *PathsSetter_SetPaths_Call {
	return &PathsSetter_SetPaths_Call{Call: _e.mock.On("SetPaths", paths)}
}

func (_c *PathsSetter_SetPaths_Call) Run(run func(paths map[string]io.Reader)) *PathsSetter_SetPaths_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]io.Reader))
	})
	return _c
}

func (_c *PathsSetter_SetPaths_Call) Return(_a0 error) *PathsSetter_SetPaths_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PathsSetter_SetPaths_Call) RunAndReturn(run func(map[string]io.Reader) error) *PathsSetter_SetPaths_Call {
	_c.Call.Return(run)
	return _c
}

// NewPathsSetter creates a new instance of PathsSetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPathsSetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *PathsSetter {
	mock := &PathsSetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
