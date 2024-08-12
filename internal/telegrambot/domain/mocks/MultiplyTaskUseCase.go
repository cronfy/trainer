// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	domain "github.com/cronfy/trainer/internal/app/domain"
	mock "github.com/stretchr/testify/mock"
)

// MultiplyTaskUseCase is an autogenerated mock type for the MultiplyTaskUseCase type
type MultiplyTaskUseCase struct {
	mock.Mock
}

type MultiplyTaskUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *MultiplyTaskUseCase) EXPECT() *MultiplyTaskUseCase_Expecter {
	return &MultiplyTaskUseCase_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields:
func (_m *MultiplyTaskUseCase) Get() domain.MultiplyTask {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 domain.MultiplyTask
	if rf, ok := ret.Get(0).(func() domain.MultiplyTask); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(domain.MultiplyTask)
	}

	return r0
}

// MultiplyTaskUseCase_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MultiplyTaskUseCase_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
func (_e *MultiplyTaskUseCase_Expecter) Get() *MultiplyTaskUseCase_Get_Call {
	return &MultiplyTaskUseCase_Get_Call{Call: _e.mock.On("Get")}
}

func (_c *MultiplyTaskUseCase_Get_Call) Run(run func()) *MultiplyTaskUseCase_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MultiplyTaskUseCase_Get_Call) Return(_a0 domain.MultiplyTask) *MultiplyTaskUseCase_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MultiplyTaskUseCase_Get_Call) RunAndReturn(run func() domain.MultiplyTask) *MultiplyTaskUseCase_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Solve provides a mock function with given fields: task, solution
func (_m *MultiplyTaskUseCase) Solve(task domain.MultiplyTask, solution int) bool {
	ret := _m.Called(task, solution)

	if len(ret) == 0 {
		panic("no return value specified for Solve")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(domain.MultiplyTask, int) bool); ok {
		r0 = rf(task, solution)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MultiplyTaskUseCase_Solve_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Solve'
type MultiplyTaskUseCase_Solve_Call struct {
	*mock.Call
}

// Solve is a helper method to define mock.On call
//   - task domain.MultiplyTask
//   - solution int
func (_e *MultiplyTaskUseCase_Expecter) Solve(task interface{}, solution interface{}) *MultiplyTaskUseCase_Solve_Call {
	return &MultiplyTaskUseCase_Solve_Call{Call: _e.mock.On("Solve", task, solution)}
}

func (_c *MultiplyTaskUseCase_Solve_Call) Run(run func(task domain.MultiplyTask, solution int)) *MultiplyTaskUseCase_Solve_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.MultiplyTask), args[1].(int))
	})
	return _c
}

func (_c *MultiplyTaskUseCase_Solve_Call) Return(_a0 bool) *MultiplyTaskUseCase_Solve_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MultiplyTaskUseCase_Solve_Call) RunAndReturn(run func(domain.MultiplyTask, int) bool) *MultiplyTaskUseCase_Solve_Call {
	_c.Call.Return(run)
	return _c
}

// NewMultiplyTaskUseCase creates a new instance of MultiplyTaskUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMultiplyTaskUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MultiplyTaskUseCase {
	mock := &MultiplyTaskUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
