// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/open-Q/user/storage/model"
	mock "github.com/stretchr/testify/mock"

	storage "github.com/open-Q/user/storage"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, user
func (_m *User) Add(ctx context.Context, user model.User) (*storage.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *storage.User
	if rf, ok := ret.Get(0).(func(context.Context, model.User) *storage.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, userID
func (_m *User) Delete(ctx context.Context, userID string) (*storage.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *storage.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *storage.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Disconnect provides a mock function with given fields: ctx
func (_m *User) Disconnect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, filter
func (_m *User) Find(ctx context.Context, filter model.UserFindFilter) ([]storage.User, error) {
	ret := _m.Called(ctx, filter)

	var r0 []storage.User
	if rf, ok := ret.Get(0).(func(context.Context, model.UserFindFilter) []storage.User); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.UserFindFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, user
func (_m *User) Update(ctx context.Context, user model.User) (*storage.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *storage.User
	if rf, ok := ret.Get(0).(func(context.Context, model.User) *storage.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
