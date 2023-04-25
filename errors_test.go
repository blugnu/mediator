package mediator

import (
	"errors"
	"fmt"
	"testing"
)

func Test_NoHandlerError(t *testing.T) {

	// ARRANGE
	request := "request"

	// ACT
	err := &NoHandlerError{"request", *new(NoResult)}

	// ASSERT
	wanted := fmt.Sprintf("no handler for '%T' request returning 'mediator.NoResult'", request)
	got := err.Error()
	if got != wanted {
		t.Errorf("wanted %q, got %q", wanted, got)
	}
}

// func Test_NoReceiverError(t *testing.T) {

// 	// ARRANGE
// 	data := "request"

// 	// ACT

// 	err := &NoReceiverError{data: data}

// 	// ASSERT

// 	wanted := fmt.Sprintf("no receiver for '%T'", data)
// 	got := err.Error()
// 	if got != wanted {
// 		t.Errorf("wanted %q, got %q", wanted, got)
// 	}
// }

func Test_InvalidHandlerError(t *testing.T) {

	// ARRANGE

	request := "request"
	handlerResult := "result"
	improperResult := true
	mock := MockCommandResult[string](handlerResult)
	defer mock.Remove()

	// ACT

	err := &ResultTypeError{handler: mock, request: request, result: true}

	// ASSERT

	wanted := fmt.Sprintf("handler for %T (%T) does not return %T", request, mock, improperResult)
	got := err.Error()
	if got != wanted {
		t.Errorf("wanted %q, got %q", wanted, got)
	}
}

func Test_ConfigurationError(t *testing.T) {
	// ARRANGE
	inner := errors.New("inner error")

	// ACT
	err := &ConfigurationError{handler: "handler", E: inner}

	// ASSERT
	wanted := "string configuration error: inner error"
	got := err.Error()
	if wanted != got {
		t.Errorf("wanted %q, got %q", wanted, got)
	}

	t.Run("unwraps the wrapped error", func(t *testing.T) {
		wanted := inner
		got := errors.Unwrap(err)
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})
}

func Test_ValidationError(t *testing.T) {
	// ARRANGE
	inner := errors.New("inner error")

	// ACT
	err := &ValidationError{inner}
	result := err.Error()

	// ASSERT
	wanted := fmt.Sprintf("validation error: %v", inner)
	got := result
	if wanted != got {
		t.Errorf("wanted %q, got %q", wanted, got)
	}

	t.Run("unwraps the wrapped error", func(t *testing.T) {
		wanted := inner
		got := errors.Unwrap(err)
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})
}
