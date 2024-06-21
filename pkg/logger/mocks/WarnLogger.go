// Code generated by mockery v2.42.1. DO NOT EDIT.

package loggermocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// WarnLogger is an autogenerated mock type for the WarnLogger type
type WarnLogger struct {
	mock.Mock
}

type WarnLogger_Expecter struct {
	mock *mock.Mock
}

func (_m *WarnLogger) EXPECT() *WarnLogger_Expecter {
	return &WarnLogger_Expecter{mock: &_m.Mock}
}

// Warn provides a mock function with given fields: msg, args
func (_m *WarnLogger) Warn(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// WarnLogger_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type WarnLogger_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *WarnLogger_Expecter) Warn(msg interface{}, args ...interface{}) *WarnLogger_Warn_Call {
	return &WarnLogger_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{msg}, args...)...)}
}

func (_c *WarnLogger_Warn_Call) Run(run func(msg string, args ...interface{})) *WarnLogger_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *WarnLogger_Warn_Call) Return() *WarnLogger_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *WarnLogger_Warn_Call) RunAndReturn(run func(string, ...interface{})) *WarnLogger_Warn_Call {
	_c.Call.Return(run)
	return _c
}

// WarnContext provides a mock function with given fields: ctx, msg, args
func (_m *WarnLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, ctx, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// WarnLogger_WarnContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WarnContext'
type WarnLogger_WarnContext_Call struct {
	*mock.Call
}

// WarnContext is a helper method to define mock.On call
//   - ctx context.Context
//   - msg string
//   - args ...interface{}
func (_e *WarnLogger_Expecter) WarnContext(ctx interface{}, msg interface{}, args ...interface{}) *WarnLogger_WarnContext_Call {
	return &WarnLogger_WarnContext_Call{Call: _e.mock.On("WarnContext",
		append([]interface{}{ctx, msg}, args...)...)}
}

func (_c *WarnLogger_WarnContext_Call) Run(run func(ctx context.Context, msg string, args ...interface{})) *WarnLogger_WarnContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *WarnLogger_WarnContext_Call) Return() *WarnLogger_WarnContext_Call {
	_c.Call.Return()
	return _c
}

func (_c *WarnLogger_WarnContext_Call) RunAndReturn(run func(context.Context, string, ...interface{})) *WarnLogger_WarnContext_Call {
	_c.Call.Return(run)
	return _c
}

// NewWarnLogger creates a new instance of WarnLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWarnLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *WarnLogger {
	mock := &WarnLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
