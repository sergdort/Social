// Code generated by mockery v2.53.3. DO NOT EDIT.

package domain

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockPostsRepository is an autogenerated mock type for the PostsRepository type
type MockPostsRepository struct {
	mock.Mock
}

type MockPostsRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPostsRepository) EXPECT() *MockPostsRepository_Expecter {
	return &MockPostsRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, post
func (_m *MockPostsRepository) Create(ctx context.Context, post *Post) error {
	ret := _m.Called(ctx, post)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPostsRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockPostsRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - post *Post
func (_e *MockPostsRepository_Expecter) Create(ctx interface{}, post interface{}) *MockPostsRepository_Create_Call {
	return &MockPostsRepository_Create_Call{Call: _e.mock.On("Create", ctx, post)}
}

func (_c *MockPostsRepository_Create_Call) Run(run func(ctx context.Context, post *Post)) *MockPostsRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*Post))
	})
	return _c
}

func (_c *MockPostsRepository_Create_Call) Return(_a0 error) *MockPostsRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPostsRepository_Create_Call) RunAndReturn(run func(context.Context, *Post) error) *MockPostsRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockPostsRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPostsRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockPostsRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockPostsRepository_Expecter) Delete(ctx interface{}, id interface{}) *MockPostsRepository_Delete_Call {
	return &MockPostsRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockPostsRepository_Delete_Call) Run(run func(ctx context.Context, id int64)) *MockPostsRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockPostsRepository_Delete_Call) Return(_a0 error) *MockPostsRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPostsRepository_Delete_Call) RunAndReturn(run func(context.Context, int64) error) *MockPostsRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockPostsRepository) GetByID(ctx context.Context, id int64) (*Post, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*Post, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *Post); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPostsRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockPostsRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockPostsRepository_Expecter) GetByID(ctx interface{}, id interface{}) *MockPostsRepository_GetByID_Call {
	return &MockPostsRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockPostsRepository_GetByID_Call) Run(run func(ctx context.Context, id int64)) *MockPostsRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockPostsRepository_GetByID_Call) Return(_a0 *Post, _a1 error) *MockPostsRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPostsRepository_GetByID_Call) RunAndReturn(run func(context.Context, int64) (*Post, error)) *MockPostsRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserFeed provides a mock function with given fields: ctx, userId, query
func (_m *MockPostsRepository) GetUserFeed(ctx context.Context, userId int64, query PaginatedFeedQuery) ([]PostWithMetadata, error) {
	ret := _m.Called(ctx, userId, query)

	if len(ret) == 0 {
		panic("no return value specified for GetUserFeed")
	}

	var r0 []PostWithMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)); ok {
		return rf(ctx, userId, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, PaginatedFeedQuery) []PostWithMetadata); ok {
		r0 = rf(ctx, userId, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PostWithMetadata)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, PaginatedFeedQuery) error); ok {
		r1 = rf(ctx, userId, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPostsRepository_GetUserFeed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserFeed'
type MockPostsRepository_GetUserFeed_Call struct {
	*mock.Call
}

// GetUserFeed is a helper method to define mock.On call
//   - ctx context.Context
//   - userId int64
//   - query PaginatedFeedQuery
func (_e *MockPostsRepository_Expecter) GetUserFeed(ctx interface{}, userId interface{}, query interface{}) *MockPostsRepository_GetUserFeed_Call {
	return &MockPostsRepository_GetUserFeed_Call{Call: _e.mock.On("GetUserFeed", ctx, userId, query)}
}

func (_c *MockPostsRepository_GetUserFeed_Call) Run(run func(ctx context.Context, userId int64, query PaginatedFeedQuery)) *MockPostsRepository_GetUserFeed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(PaginatedFeedQuery))
	})
	return _c
}

func (_c *MockPostsRepository_GetUserFeed_Call) Return(_a0 []PostWithMetadata, _a1 error) *MockPostsRepository_GetUserFeed_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPostsRepository_GetUserFeed_Call) RunAndReturn(run func(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)) *MockPostsRepository_GetUserFeed_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, post
func (_m *MockPostsRepository) Update(ctx context.Context, post *Post) error {
	ret := _m.Called(ctx, post)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Post) error); ok {
		r0 = rf(ctx, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPostsRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockPostsRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - post *Post
func (_e *MockPostsRepository_Expecter) Update(ctx interface{}, post interface{}) *MockPostsRepository_Update_Call {
	return &MockPostsRepository_Update_Call{Call: _e.mock.On("Update", ctx, post)}
}

func (_c *MockPostsRepository_Update_Call) Run(run func(ctx context.Context, post *Post)) *MockPostsRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*Post))
	})
	return _c
}

func (_c *MockPostsRepository_Update_Call) Return(_a0 error) *MockPostsRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPostsRepository_Update_Call) RunAndReturn(run func(context.Context, *Post) error) *MockPostsRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPostsRepository creates a new instance of MockPostsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPostsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPostsRepository {
	mock := &MockPostsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
