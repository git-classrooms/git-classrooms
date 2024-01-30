// Code generated by mockery v2.40.1. DO NOT EDIT.

package api

import (
	fiber "github.com/gofiber/fiber/v2"
	mock "github.com/stretchr/testify/mock"
)

// MockController is an autogenerated mock type for the Controller type
type MockController struct {
	mock.Mock
}

type MockController_Expecter struct {
	mock *mock.Mock
}

func (_m *MockController) EXPECT() *MockController_Expecter {
	return &MockController_Expecter{mock: &_m.Mock}
}

// CreateAssignment provides a mock function with given fields: c
func (_m *MockController) CreateAssignment(c *fiber.Ctx) error {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for CreateAssignment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*fiber.Ctx) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockController_CreateAssignment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAssignment'
type MockController_CreateAssignment_Call struct {
	*mock.Call
}

// CreateAssignment is a helper method to define mock.On call
//   - c *fiber.Ctx
func (_e *MockController_Expecter) CreateAssignment(c interface{}) *MockController_CreateAssignment_Call {
	return &MockController_CreateAssignment_Call{Call: _e.mock.On("CreateAssignment", c)}
}

func (_c *MockController_CreateAssignment_Call) Run(run func(c *fiber.Ctx)) *MockController_CreateAssignment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*fiber.Ctx))
	})
	return _c
}

func (_c *MockController_CreateAssignment_Call) Return(_a0 error) *MockController_CreateAssignment_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockController_CreateAssignment_Call) RunAndReturn(run func(*fiber.Ctx) error) *MockController_CreateAssignment_Call {
	_c.Call.Return(run)
	return _c
}

// CreateClassroom provides a mock function with given fields: c
func (_m *MockController) CreateClassroom(c *fiber.Ctx) error {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for CreateClassroom")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*fiber.Ctx) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockController_CreateClassroom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateClassroom'
type MockController_CreateClassroom_Call struct {
	*mock.Call
}

// CreateClassroom is a helper method to define mock.On call
//   - c *fiber.Ctx
func (_e *MockController_Expecter) CreateClassroom(c interface{}) *MockController_CreateClassroom_Call {
	return &MockController_CreateClassroom_Call{Call: _e.mock.On("CreateClassroom", c)}
}

func (_c *MockController_CreateClassroom_Call) Run(run func(c *fiber.Ctx)) *MockController_CreateClassroom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*fiber.Ctx))
	})
	return _c
}

func (_c *MockController_CreateClassroom_Call) Return(_a0 error) *MockController_CreateClassroom_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockController_CreateClassroom_Call) RunAndReturn(run func(*fiber.Ctx) error) *MockController_CreateClassroom_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockController creates a new instance of MockController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockController {
	mock := &MockController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
