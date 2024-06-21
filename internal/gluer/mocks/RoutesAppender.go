// Code generated by mockery v2.42.1. DO NOT EDIT.

package gluermocks

import (
	route "github.com/amidgo/swaglue/internal/route"
	mock "github.com/stretchr/testify/mock"
)

// RoutesAppender is an autogenerated mock type for the RoutesAppender type
type RoutesAppender struct {
	mock.Mock
}

type RoutesAppender_Expecter struct {
	mock *mock.Mock
}

func (_m *RoutesAppender) EXPECT() *RoutesAppender_Expecter {
	return &RoutesAppender_Expecter{mock: &_m.Mock}
}

// AppendRoutes provides a mock function with given fields: routes
func (_m *RoutesAppender) AppendRoutes(routes []route.Route) error {
	ret := _m.Called(routes)

	if len(ret) == 0 {
		panic("no return value specified for AppendRoutes")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]route.Route) error); ok {
		r0 = rf(routes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoutesAppender_AppendRoutes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendRoutes'
type RoutesAppender_AppendRoutes_Call struct {
	*mock.Call
}

// AppendRoutes is a helper method to define mock.On call
//   - routes []route.Route
func (_e *RoutesAppender_Expecter) AppendRoutes(routes interface{}) *RoutesAppender_AppendRoutes_Call {
	return &RoutesAppender_AppendRoutes_Call{Call: _e.mock.On("AppendRoutes", routes)}
}

func (_c *RoutesAppender_AppendRoutes_Call) Run(run func(routes []route.Route)) *RoutesAppender_AppendRoutes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]route.Route))
	})
	return _c
}

func (_c *RoutesAppender_AppendRoutes_Call) Return(_a0 error) *RoutesAppender_AppendRoutes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoutesAppender_AppendRoutes_Call) RunAndReturn(run func([]route.Route) error) *RoutesAppender_AppendRoutes_Call {
	_c.Call.Return(run)
	return _c
}

// NewRoutesAppender creates a new instance of RoutesAppender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRoutesAppender(t interface {
	mock.TestingT
	Cleanup(func())
}) *RoutesAppender {
	mock := &RoutesAppender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
