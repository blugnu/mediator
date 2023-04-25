package mediator

import "context"

// CommandHandlerFunc[TRequest, TResult] is the function signature of a
// command handler.
type CommandHandlerFunc[TRequest any, TResult any] func(context.Context, TRequest) (TResult, error)

// ValidatorFunc[TInput] is the function signature of a validator.
type ValidatorFunc[TRequest any] func(context.Context, TRequest) error

// ValidatorFunc[TInput] is the function signature of a validator.
type ConfigurationCheckerFunc[TRequest any] func(context.Context) error

// mockCommand[TRequest any, TResult any] is a mock command handler that can be
// used in tests to verify that a command is called with the expected parameters
// and/or to return a specified result or error.
type mockcommand[TRequest any, TResult any] struct {
	requests           []TRequest
	checkConfiguration ConfigurationCheckerFunc[TRequest]
	validate           ValidatorFunc[TRequest]
	execute            CommandHandlerFunc[TRequest, TResult]
	unregister         func()
}

// CheckConfiguration satisfies the ConfigurationChecker interface
func (mock *mockcommand[TRequest, TResult]) CheckConfiguration(ctx context.Context) error {
	if mock.checkConfiguration != nil {
		return mock.checkConfiguration(ctx)
	}
	return nil
}

// Validate satisfies the Validator interface
func (mock *mockcommand[TRequest, TResult]) Validate(ctx context.Context, rq TRequest) error {
	mock.requests = append(mock.requests, rq)
	if mock.validate != nil {
		return mock.validate(ctx, rq)
	}
	return nil
}

// Execute satisfies the CommandHandler interface
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

func (mock *mockcommand[TRequest, TResult]) Remove() {
	mock.unregister()
}

// mockCommandHandler registers a mock handler for the specified request and result
// type using the specified handler function.
func mockCommandHandler[TRequest any, TResult any](cfg ConfigurationCheckerFunc[TRequest], v ValidatorFunc[TRequest], cmd CommandHandlerFunc[TRequest, TResult]) *mockcommand[TRequest, TResult] {
	mock := &mockcommand[TRequest, TResult]{
		checkConfiguration: cfg,
		validate:           v,
		execute:            cmd,
	}
	mock.unregister = RegisterCommand[TRequest, TResult](mock)
	return mock
}

// MockCommand registers a mock handler for the specified request and result type
// modelling successful execution of the command returning a zero-value result
// and no error.
func MockCommand[TRequest any, TResult any]() *mockcommand[TRequest, TResult] {
	return MockCommandResult[TRequest](*new(TResult))
}

// MockCommandError registers a mock handler for the specified request and result
// type modelling unsuccessful execution of the command returning the specified
// error.
func MockCommandError[TRequest any, TResult any](err error) *mockcommand[TRequest, TResult] {
	// return mockCommandHandler(nil, func(ctx context.Context, rq TRequest) (TResult, error) { return *new(TResult), err })
	return mockCommandHandler(nil, nil, func(ctx context.Context, rq TRequest) (TResult, error) { return *new(TResult), err })
}

// MockCommandResult registers a mock handler for the specified request and result
// type modelling successful execution of the command returning the specified
// result and no error.
func MockCommandResult[TRequest any, TResult any](result TResult) *mockcommand[TRequest, TResult] {
	return mockCommandHandler(nil, nil, func(ctx context.Context, rq TRequest) (TResult, error) { return result, nil })
}

// MockCommandConfigurationError registers a mock handler for the specified request
// and result type modelling an unsuccessful configuration check of the command.
// The specified error will be returned by the configuration checker of the mocked
// command (and will therefore be wrapped in a ConfigurationError).
func MockCommandConfigurationError[TRequest any, TResult any](err error) *mockcommand[TRequest, TResult] {
	return mockCommandHandler[TRequest, TResult](func(context.Context) error { return err }, nil, nil)
}

// MockCommandValidationError registers a mock handler for the specified request
// and result type modelling unsuccessful validation of the command.  The specified
// error will be returned by the validator of the mocked command (and will therefore
// be wrapped in a ValidationError).
func MockCommandValidationError[TRequest any, TResult any](err error) *mockcommand[TRequest, TResult] {
	return mockCommandHandler[TRequest, TResult](nil, func(context.Context, TRequest) error { return err }, nil)
}
