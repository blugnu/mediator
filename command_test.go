package mediator

import (
	"context"
	"errors"
	"testing"
)

// MockCommandConfigurationChecker registers a mock handler for the specified request
// and result type using the specified configuration checker function.
func MockCommandConfigurationChecker[TRequest any, TResult any](qry CommandHandlerFunc[TRequest, TResult], cfg ConfigurationCheckerFunc[TRequest]) *mockcommand[TRequest, TResult] {
	return mockCommandHandler(cfg, nil, func(context.Context, TRequest) (TResult, error) { return *new(TResult), nil })
}

// MockCommandValidator registers a mock handler for the specified request and result
// type using the specified validator function.
func MockCommandValidator[TRequest any, TResult any](qry CommandHandlerFunc[TRequest, TResult], validator ValidatorFunc[TRequest]) *mockcommand[TRequest, TResult] {
	return mockCommandHandler(nil, validator, func(context.Context, TRequest) (TResult, error) { return *new(TResult), nil })
}

func TestThatTheRegistrationInterfaceRemovesTheHandler(t *testing.T) {

	if len(commandHandlers) > 0 {
		t.Fatal("invalid test: one or more command handlers are already registered")
	}

	// ARRANGE
	mock := MockCommandResult[string]("")

	wanted := 1
	got := len(commandHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}

	// ACT
	mock.Remove()

	// ASSERT
	wanted = 0
	got = len(commandHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
}

func TestThatRegisterHandlerPanicsWhenHandlerIsAlreadyRegisteredForAType(t *testing.T) {
	// ARRANGE
	defer func() { // panic tests must be deferred
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()

	// Register a handler and remove it when done
	mock := MockCommandResult[string]("result")
	defer mock.Remove()

	// ACT - attempt to register another handler for the same request type
	MockCommandResult[string]("other")

	// ASSERT (deferred, see above)
}

func TestThatHandlerReturnsExpectedErrorWhenHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT
	_, err := Execute(context.Background(), "request", new(NoResult))

	// ASSERT
	wanted := &NoHandlerError{"request", *new(NoResult)}
	got := err
	if !errors.Is(got, wanted) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}

func TestThatHandlerReturnsExpectedErrorWhenHandlerResultIsWrongType(t *testing.T) {
	// ARRANGE

	// Register a handler returning a string
	mock := MockCommandResult[string]("string response")
	defer mock.Remove()

	// ACT

	// Request a Handler returning a bool
	_, err := Execute(context.Background(), "request", new(NoResult))

	// ASSERT

	if _, ok := err.(*ResultTypeError); !ok {
		t.Errorf("wanted *mediator.ErrInvalidHandler, got %T", err)
	}
}

func TestThatConfigurationCheckerErrorIsReturned(t *testing.T) {
	// ARRANGE
	mock := MockCommandConfigurationError[string, NoResult](errors.New("error"))
	defer mock.Remove()

	// ACT
	_, err := Execute(context.Background(), "request", new(NoResult))

	// ASSERT
	wanted := ConfigurationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}
}

func TestThatValidatorErrorIsReturned(t *testing.T) {
	// ARRANGE
	mock := MockCommandValidationError[string, NoResult](errors.New("error"))
	defer mock.Remove()

	// ACT
	_, err := Execute(context.Background(), "request", new(NoResult))

	// ASSERT
	wanted := ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}
}

func TestThatMediatorDoesNotWrapConfigurationCheckerError(t *testing.T) {
	// ARRANGE
	mock := MockCommandConfigurationError[string, NoResult](ConfigurationError{})
	defer mock.Remove()

	// ACT
	_, err := Execute(context.Background(), "request", new(NoResult))

	// ASSERT
	var wanted = ConfigurationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("\nwanted %T\ngot    %T (%[2]q)", wanted, err)
	}

	got := wanted.E
	if errors.As(got, &wanted) {
		t.Errorf("got %T wrapping %[1]T unnecessarily (%[1]v)", err)
	}
}

func TestThatHandlerValidatorDoesNotWrapValidationError(t *testing.T) {
	// ARRANGE
	mock := MockCommandValidationError[string, NoResult](ValidationError{})
	defer mock.Remove()

	// ACT
	_, err := Execute(context.Background(), "request", new(NoResult))

	// ASSERT
	var wanted = ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("\nwanted %T\ngot    %T (%[2]q)", wanted, err)
	}

	got := wanted.E
	if errors.As(got, &wanted) {
		t.Errorf("got %T wrapping %[1]T unnecessarily (%[1]v)", err)
	}
}

func TestThatHandlerResultIsReturnedInResultParam(t *testing.T) {
	// ARRANGE
	wanted := "result"
	mock := MockCommandResult[string](wanted)
	defer mock.Remove()

	// ACT
	result, err := Execute(context.Background(), "request", new(string))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// ASSERT
	got := result
	if wanted != got {
		t.Errorf("wanted %q, got %q", wanted, got)
	}
}
