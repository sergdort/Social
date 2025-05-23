// Code generated by mockery v2.53.3. DO NOT EDIT.

package domain

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCommentsRepository is an autogenerated mock type for the CommentsRepository type
type MockCommentsRepository struct {
	mock.Mock
}

type MockCommentsRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCommentsRepository) EXPECT() *MockCommentsRepository_Expecter {
	return &MockCommentsRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, comment
func (_m *MockCommentsRepository) Create(ctx context.Context, comment *Comment) error {
	ret := _m.Called(ctx, comment)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Comment) error); ok {
		r0 = rf(ctx, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCommentsRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockCommentsRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - comment *Comment
func (_e *MockCommentsRepository_Expecter) Create(ctx interface{}, comment interface{}) *MockCommentsRepository_Create_Call {
	return &MockCommentsRepository_Create_Call{Call: _e.mock.On("Create", ctx, comment)}
}

func (_c *MockCommentsRepository_Create_Call) Run(run func(ctx context.Context, comment *Comment)) *MockCommentsRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*Comment))
	})
	return _c
}

func (_c *MockCommentsRepository_Create_Call) Return(_a0 error) *MockCommentsRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommentsRepository_Create_Call) RunAndReturn(run func(context.Context, *Comment) error) *MockCommentsRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByPostID provides a mock function with given fields: ctx, postID
func (_m *MockCommentsRepository) GetAllByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	ret := _m.Called(ctx, postID)

	if len(ret) == 0 {
		panic("no return value specified for GetAllByPostID")
	}

	var r0 []Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]Comment, error)); ok {
		return rf(ctx, postID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []Comment); ok {
		r0 = rf(ctx, postID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCommentsRepository_GetAllByPostID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByPostID'
type MockCommentsRepository_GetAllByPostID_Call struct {
	*mock.Call
}

// GetAllByPostID is a helper method to define mock.On call
//   - ctx context.Context
//   - postID int64
func (_e *MockCommentsRepository_Expecter) GetAllByPostID(ctx interface{}, postID interface{}) *MockCommentsRepository_GetAllByPostID_Call {
	return &MockCommentsRepository_GetAllByPostID_Call{Call: _e.mock.On("GetAllByPostID", ctx, postID)}
}

func (_c *MockCommentsRepository_GetAllByPostID_Call) Run(run func(ctx context.Context, postID int64)) *MockCommentsRepository_GetAllByPostID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockCommentsRepository_GetAllByPostID_Call) Return(_a0 []Comment, _a1 error) *MockCommentsRepository_GetAllByPostID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCommentsRepository_GetAllByPostID_Call) RunAndReturn(run func(context.Context, int64) ([]Comment, error)) *MockCommentsRepository_GetAllByPostID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCommentsRepository creates a new instance of MockCommentsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCommentsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCommentsRepository {
	mock := &MockCommentsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
