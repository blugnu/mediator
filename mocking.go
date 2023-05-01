package mediator

import (
	"context"
	"fmt"
)

// executeFunc[TRequest, TResult] is the signature of a function that implements
// the CommandExecutor[TRequest, TResult] interface.
type executeFunc[TRequest any, TResult any] func(context.Context, TRequest) (TResult, error)

// validateFunc[TInput] is the signature of a function that implements
// the Validator[TRequest] interface.
type validateFunc[TRequest any] func(context.Context, TRequest) error

// mockCommand[TRequest any, TResult any] is a mock command that can be
// used in tests to verify that a command is called with the expected parameters
// and/or to return a specified result or error.
type mockcommand[TRequest any, TResult any] struct {
	requests   []TRequest
	validate   validateFunc[TRequest]
	execute    executeFunc[TRequest, TResult]
	unregister func()
}

// Validate satisfies the Validator interface
func (mock *mockcommand[TRequest, TResult]) Validate(ctx context.Context, rq TRequest) error {
	mock.requests = append(mock.requests, rq)
	if mock.validate != nil {
		return mock.validate(ctx, rq)
	}
	return nil
}

// Execute satisfies the CommandExecutor interface
func (mock *mockcommand[TRequest, TResult]) Execute(ctx context.Context, rq TRequest) (TResult, error) {
	return mock.execute(ctx, rq)
}

// NumRequests returns the number of times the mock was called.
func (mock *mockcommand[TRequest, TResult]) NumRequests() int {
	return len(mock.requests)
}

// Requests returns a copy of the slice of requests received by the mock.
func (mock *mockcommand[TRequest, TResult]) Requests() []TRequest {
	return append([]TRequest{}, mock.requests...)
}

// WasCalled returns true if the mock was called at least once.
func (mock *mockcommand[TRequest, TResult]) WasCalled() bool {
	return len(mock.requests) > 0
}

// WasNotCalled returns true if the mock was not called.
func (mock *mockcommand[TRequest, TResult]) WasNotCalled() bool {
	return len(mock.requests) == 0
}

func (mock *mockcommand[TRequest, TResult]) Unregister() {
	mock.unregister()
}

// registerMockCommand registers a mock command for the specified request and result
// type using the specified command functions for validating requests and executing
// the command.
func registerMockCommand[TRequest any, TResult any](val validateFunc[TRequest], cmd executeFunc[TRequest, TResult]) *mockcommand[TRequest, TResult] {
	mock := &mockcommand[TRequest, TResult]{
		validate: val,
		execute:  cmd,
	}
	mock.unregister, _ = register(context.Background(), *new(TRequest), mock)
	return mock
}

// MockCommand registers a mock command for the specified request and result type
// modelling successful execution of the command returning a zero-value result
// and no error.
func MockCommand[TRequest any, TResult any]() *mockcommand[TRequest, TResult] {
	return MockCommandResult[TRequest](*new(TResult))
}

// MockCommandError registers a mock command for the specified request and result
// type modelling a failed execution of the command, returning the specified
// error and a zero-value result.
func MockCommandError[TRequest any, TResult any](err error) *mockcommand[TRequest, TResult] {
	return registerMockCommand(nil, func(ctx context.Context, rq TRequest) (TResult, error) { return *new(TResult), err })
}

// MockCommandResult registers a mock command for the specified request and result
// type modelling successful execution of the command returning the specified
// result and a nil error.
func MockCommandResult[TRequest any, TResult any](result TResult) *mockcommand[TRequest, TResult] {
	return registerMockCommand(nil, func(ctx context.Context, rq TRequest) (TResult, error) { return result, nil })
}

// MockCommandValidationError registers a mock command for the specified request
// and result type modelling a failed validation of the request.
//
// The specified error will be returned by the validator of the mocked command (and
// will therefore be wrapped in a ValidationError).
func MockCommandValidationError[TRequest any, TResult any](err error) *mockcommand[TRequest, TResult] {
	return registerMockCommand[TRequest, TResult](func(context.Context, TRequest) error { return err }, nil)
}

// RegisterMockCommand registers a custom mock command, returning a function to unregister
// the mock when no longer required:
//
//	unreg := mediator.RegisterMockCommand[MyRequest, MyResult](myMock)
//	defer unreg()
//
// Custom mocks are useful when you want to mock a command that has a custom validator
// or executor.  If a custom mock implements CheckConfiguration, it must not return any
// error from the configuration check.
//
// If a custom mock fails configuration checks the registration will panic.
func RegisterMockCommand[TRequest any, TResult any](ctx context.Context, mock CommandHandler[TRequest, TResult]) func() {
	fn, err := register(ctx, *new(TRequest), mock)
	if err != nil {
		panic(fmt.Sprintf("%T returned an error from CheckConfiguration(): %v\nCustom command mocks must not implement failing configuration checks", mock, err))
	}
	return fn
}
