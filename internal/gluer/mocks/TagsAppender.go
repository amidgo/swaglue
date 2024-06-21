// Code generated by mockery v2.42.1. DO NOT EDIT.

package gluermocks

import (
	model "github.com/amidgo/swaglue/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// TagsAppender is an autogenerated mock type for the TagsAppender type
type TagsAppender struct {
	mock.Mock
}

type TagsAppender_Expecter struct {
	mock *mock.Mock
}

func (_m *TagsAppender) EXPECT() *TagsAppender_Expecter {
	return &TagsAppender_Expecter{mock: &_m.Mock}
}

// AppendTags provides a mock function with given fields: items
func (_m *TagsAppender) AppendTags(items []model.Item) error {
	ret := _m.Called(items)

	if len(ret) == 0 {
		panic("no return value specified for AppendTags")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]model.Item) error); ok {
		r0 = rf(items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TagsAppender_AppendTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendTags'
type TagsAppender_AppendTags_Call struct {
	*mock.Call
}

// AppendTags is a helper method to define mock.On call
//   - items []model.Item
func (_e *TagsAppender_Expecter) AppendTags(items interface{}) *TagsAppender_AppendTags_Call {
	return &TagsAppender_AppendTags_Call{Call: _e.mock.On("AppendTags", items)}
}

func (_c *TagsAppender_AppendTags_Call) Run(run func(items []model.Item)) *TagsAppender_AppendTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]model.Item))
	})
	return _c
}

func (_c *TagsAppender_AppendTags_Call) Return(_a0 error) *TagsAppender_AppendTags_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TagsAppender_AppendTags_Call) RunAndReturn(run func([]model.Item) error) *TagsAppender_AppendTags_Call {
	_c.Call.Return(run)
	return _c
}

// NewTagsAppender creates a new instance of TagsAppender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTagsAppender(t interface {
	mock.TestingT
	Cleanup(func())
}) *TagsAppender {
	mock := &TagsAppender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
