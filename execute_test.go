package mediator

import (
	"context"
	"errors"
	"testing"
)

func TestThatUnregisterRemovesTheCommand(t *testing.T) {

	if len(commands) > 0 {
		t.Fatal("invalid test: one or more commands are already registered")
	}

	// ARRANGE
	mock := MockCommandResult[string]("")

	wanted := 1
	got := len(commands)
	if wanted != got {
		t.Errorf("wanted %d registered commands, got %d", wanted, got)
	}

	// ACT
	mock.Unregister()

	// ASSERT
	wanted = 0
	got = len(commands)
	if wanted != got {
		t.Errorf("wanted %d registered commands, got %d", wanted, got)
	}
}

func TestThatHandlerReturnsExpectedErrorWhenHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT
	_, err := Execute(context.Background(), "request", NoResult)

	// ASSERT
	wanted := &NoCommandForRequestTypeError{"request"}
	got := err
	if !errors.Is(got, wanted) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}

func TestThatHandlerReturnsExpectedErrorWhenHandlerResultIsWrongType(t *testing.T) {
	// ARRANGE - register a command returning a string
	mock := MockCommand[string, string]()
	defer mock.Unregister()

	// ACT - expect a command returning NoResult
	_, err := Execute(context.Background(), "request", NoResult)

	// ASSERT
	wanted := &ResultTypeError{mock, *new(NoResultType)}
	got := err
	if !errors.Is(got, wanted) {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
	}
}

func TestThatValidatorErrorIsReturned(t *testing.T) {
	// ARRANGE
	mock := MockCommandValidationError[string, NoResultType](errors.New("error"))
	defer mock.Unregister()

	// ACT
	_, err := Execute(context.Background(), "request", NoResult)

	// ASSERT
	wanted := ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}
}

func TestThatMediatorDoesNotWrapValidationError(t *testing.T) {
	rawerr := errors.New("error")
	testcases := []struct {
		name string
		error
	}{
		{name: "wrapped by value", error: ValidationError{E: rawerr}},
		{name: "wrapped by reference", error: &ValidationError{E: rawerr}},
		{name: "returned unwrapped", error: rawerr},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			mock := MockCommandValidationError[string, NoResultType](tc.error)
			defer mock.Unregister()

			// ACT
			_, err := Execute(context.Background(), "request", NoResult)

			// ASSERT
			var wanted = ValidationError{}
			t.Run("returns expected type", func(t *testing.T) {
				if !errors.As(err, &wanted) {
					t.Errorf("\nwanted %T\ngot    %T (%[2]q)", wanted, err)
				}
			})

			t.Run("is not wrapped", func(t *testing.T) {
				got := wanted.E
				if errors.As(got, &wanted) {
					t.Errorf("got %T wrapping %[1]T unnecessarily (%[1]v)", err)
				}
			})
		})
	}
}

func TestThatHandlerResultIsReturnedInResultParam(t *testing.T) {
	// ARRANGE
	wanted := "result"
	mock := MockCommandResult[string](wanted)
	defer mock.Unregister()

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
