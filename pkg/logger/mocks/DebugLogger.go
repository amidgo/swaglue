// Code generated by mockery v2.42.1. DO NOT EDIT.

package loggermocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DebugLogger is an autogenerated mock type for the DebugLogger type
type DebugLogger struct {
	mock.Mock
}

type DebugLogger_Expecter struct {
	mock *mock.Mock
}

func (_m *DebugLogger) EXPECT() *DebugLogger_Expecter {
	return &DebugLogger_Expecter{mock: &_m.Mock}
}

// Debug provides a mock function with given fields: msg, args
func (_m *DebugLogger) Debug(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// DebugLogger_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type DebugLogger_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *DebugLogger_Expecter) Debug(msg interface{}, args ...interface{}) *DebugLogger_Debug_Call {
	return &DebugLogger_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{msg}, args...)...)}
}

func (_c *DebugLogger_Debug_Call) Run(run func(msg string, args ...interface{})) *DebugLogger_Debug_Call {
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

func (_c *DebugLogger_Debug_Call) Return() *DebugLogger_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *DebugLogger_Debug_Call) RunAndReturn(run func(string, ...interface{})) *DebugLogger_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// DebugContext provides a mock function with given fields: ctx, msg, args
func (_m *DebugLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, ctx, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// DebugLogger_DebugContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DebugContext'
type DebugLogger_DebugContext_Call struct {
	*mock.Call
}

// DebugContext is a helper method to define mock.On call
//   - ctx context.Context
//   - msg string
//   - args ...interface{}
func (_e *DebugLogger_Expecter) DebugContext(ctx interface{}, msg interface{}, args ...interface{}) *DebugLogger_DebugContext_Call {
	return &DebugLogger_DebugContext_Call{Call: _e.mock.On("DebugContext",
		append([]interface{}{ctx, msg}, args...)...)}
}

func (_c *DebugLogger_DebugContext_Call) Run(run func(ctx context.Context, msg string, args ...interface{})) *DebugLogger_DebugContext_Call {
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

func (_c *DebugLogger_DebugContext_Call) Return() *DebugLogger_DebugContext_Call {
	_c.Call.Return()
	return _c
}

func (_c *DebugLogger_DebugContext_Call) RunAndReturn(run func(context.Context, string, ...interface{})) *DebugLogger_DebugContext_Call {
	_c.Call.Return(run)
	return _c
}

// NewDebugLogger creates a new instance of DebugLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDebugLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *DebugLogger {
	mock := &DebugLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
