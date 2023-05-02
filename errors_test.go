package mediator

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// errorstestcmd is a command used for testing errors.
type errorstestcmd struct{}

func (errorstestcmd) Execute(context.Context, int) (NoResultType, error) { return nil, nil }

func Test_CommandAlreadyRegisteredError(t *testing.T) {
	t.Run("Error()", func(t *testing.T) {
		testcases := []struct {
			request any
			result  string
		}{
			{request: "", result: "mediator.errorstestcmd already registered for requests of type: string"},
			{request: true, result: "mediator.errorstestcmd already registered for requests of type: bool"},
			{request: 42, result: "mediator.errorstestcmd already registered for requests of type: int"},
		}
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("%T", tc.request), func(t *testing.T) {
				// ARRANGE
				sut := &CommandAlreadyRegisteredError{command: errorstestcmd{}, request: tc.request}

				// ACT
				s := sut.Error()

				// ASSERT
				wanted := tc.result
				got := s
				if wanted != got {
					t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
				}
			})
		}
	})

	t.Run("Is(target)", func(t *testing.T) {
		// ARRANGE
		sut := CommandAlreadyRegisteredError{command: errorstestcmd{}, request: ""}

		testcases := []struct {
			target error
			result bool
		}{
			// same error and request type
			{target: CommandAlreadyRegisteredError{request: ""}, result: true},
			{target: &CommandAlreadyRegisteredError{request: ""}, result: true},
			// same error but different request type
			{target: CommandAlreadyRegisteredError{request: 0}, result: false},
			{target: &CommandAlreadyRegisteredError{request: 0}, result: false},
			// different error
			{target: errors.New("other error"), result: false},
		}
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("target = %T", tc.target), func(t *testing.T) {
				// ACT
				got := sut.Is(tc.target)

				// ASSERT
				wanted := tc.result
				if wanted != got {
					t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
				}
			})
		}
	})
}

func Test_NoCommandForRequestTypeError(t *testing.T) {
	t.Run("Error()", func(t *testing.T) {
		testcases := []struct {
			request any
			result  string
		}{
			{request: "", result: "no command registered for requests of type: string"},
			{request: true, result: "no command registered for requests of type: bool"},
			{request: 42, result: "no command registered for requests of type: int"},
		}
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("%T", tc.request), func(t *testing.T) {
				// ARRANGE
				sut := &NoCommandForRequestTypeError{tc.request}

				// ACT
				s := sut.Error()

				// ASSERT
				wanted := tc.result
				got := s
				if got != wanted {
					t.Errorf("wanted %q, got %q", wanted, got)
				}
			})
		}
	})

	t.Run("Is(target)", func(t *testing.T) {
		// ARRANGE
		sut := NoCommandForRequestTypeError{request: ""}

		testcases := []struct {
			target error
			result bool
		}{
			// same error and request type
			{target: NoCommandForRequestTypeError{request: "request"}, result: true},
			{target: &NoCommandForRequestTypeError{request: "request"}, result: true},
			// same error but different request type
			{target: NoCommandForRequestTypeError{request: 42}, result: false},
			{target: &NoCommandForRequestTypeError{request: 42}, result: false},
			// different error
			{target: errors.New("other error"), result: false},
		}
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("target = %T", tc.target), func(t *testing.T) {
				// ACT
				got := sut.Is(tc.target)

				// ASSERT
				wanted := tc.result
				if wanted != got {
					t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
				}
			})
		}
	})
}

func Test_ResultTypeError(t *testing.T) {
	// ARRANGE
	mock := MockCommandResult[string]("command result")
	defer mock.Unregister()

	// ACT
	sut := &ResultTypeError{command: mock, result: "command result"}

	t.Run("Error()", func(t *testing.T) {
		// ASSERT
		wanted := fmt.Sprintf("%T does not return string", mock)
		got := sut.Error()
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("Is(target)", func(t *testing.T) {
		testcases := []struct {
			target error
			result bool
		}{
			// same error and request type
			{target: ResultTypeError{command: mock, result: ""}, result: true},
			{target: &ResultTypeError{command: mock, result: ""}, result: true},
			// same error but different request type
			{target: ResultTypeError{command: mock, result: 0}, result: false},
			{target: &ResultTypeError{command: mock, result: 0}, result: false},
			// different error
			{target: errors.New("other error"), result: false},
		}
		for _, tc := range testcases {
			t.Run(fmt.Sprintf("target = %T", tc.target), func(t *testing.T) {
				// ACT
				got := sut.Is(tc.target)

				// ASSERT
				wanted := tc.result
				if wanted != got {
					t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
				}
			})
		}
	})
}

func Test_ValidationError(t *testing.T) {
	// ARRANGE
	e := errors.New("inner error")
	h := MockCommand[string, NoResultType]()
	defer h.Unregister()

	// ACT
	sut := &ValidationError{E: e}
	result := sut.Error()

	// ASSERT
	t.Run("Error()", func(t *testing.T) {
		wanted := fmt.Sprintf("request validation error: %v", e)
		got := result
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("Unwrap()", func(t *testing.T) {
		wanted := e
		got := sut.Unwrap()
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}
