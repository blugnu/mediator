package mediator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMockHandler(t *testing.T) {
	if len(commandHandlers) > 0 {
		t.Fatal("invalid test: one or more commandHandlers are already registered")
	}

	mock := MockCommand[string, NoResult]()
	defer mock.Remove()

	t.Run("registers the handler", func(t *testing.T) {
		wanted := 1
		got := len(commandHandlers)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	// ACT
	_, err := Execute(context.Background(), "test", new(NoResult))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("captures number of handled requests", func(t *testing.T) {
		wanted := 1
		got := mock.NumRequests()
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	t.Run("returns a copy of handled requests", func(t *testing.T) {
		requests := mock.Requests()

		if reflect.ValueOf(requests).UnsafePointer() == reflect.ValueOf(mock.requests).UnsafePointer() {
			t.Error("got same slice")
		}

		if !reflect.DeepEqual(requests, mock.requests) {
			t.Errorf("wanted %v, got %v", mock.requests, requests)
		}
	})

	t.Run("captures that a handler was called", func(t *testing.T) {
		wantedWc := true
		wantedWnc := false
		gotWc := mock.WasCalled()
		gotWnc := mock.WasNotCalled()

		if wantedWc != gotWc || wantedWnc != gotWnc {
			t.Errorf("called / not called: wanted %v / %v, got %v / %v", wantedWc, wantedWnc, gotWc, gotWnc)
		}
	})

	t.Run("captures that a handler was not called", func(t *testing.T) {
		mock := MockCommand[int, int]()
		defer mock.Remove()

		wantedWc := false
		wantedWnc := true
		gotWc := mock.WasCalled()
		gotWnc := mock.WasNotCalled()

		if wantedWc != gotWc || wantedWnc != gotWnc {
			t.Errorf("called / not called: wanted %v / %v, got %v / %v", wantedWc, wantedWnc, gotWc, gotWnc)
		}
	})
}

func TestMockCommandError(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	herr := errors.New("handler error")
	mock := MockCommandError[string, NoResult](herr)
	defer mock.Remove()

	// ACT
	_, err := Execute(ctx, "test", new(NoResult))

	// ASSERT
	t.Run("returns mocked error", func(t *testing.T) {
		wanted := herr
		got := err
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}

func TestMockCommandResult(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mock := MockCommandResult[string](42)
	defer mock.Remove()

	// ACT
	result, err := Execute(ctx, "test", new(int))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ASSERT
	t.Run("returns mocked error", func(t *testing.T) {
		wanted := 42
		got := result
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}

func TestMockCommandConfigurationError(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	cfgerr := errors.New("configuration error")
	mock := MockCommandConfigurationError[string, NoResult](cfgerr)
	defer mock.Remove()

	// ACT
	_, err := Execute(ctx, "test", new(NoResult))

	// ASSERT
	t.Run("does not record request", func(t *testing.T) {
		wanted := true
		got := mock.NumRequests() == 0
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("returns expected error", func(t *testing.T) {
		wanted := cfgerr
		got := err
		if !errors.Is(got, wanted) {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}

		t.Run("wrapped in ConfigurationError", func(t *testing.T) {
			wanted := ConfigurationError{}
			got := err
			if !errors.As(got, &wanted) {
				t.Errorf("\nwanted %T\ngot    %T", wanted, got)
			}
		})
	})
}

func TestMockCommandValidationError(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	verr := errors.New("validation error")
	mock := MockCommandValidationError[string, NoResult](verr)
	defer mock.Remove()

	// ACT
	_, err := Execute(ctx, "test", new(NoResult))

	// ASSERT
	t.Run("records request", func(t *testing.T) {
		wanted := true
		got := mock.NumRequests() == 1
		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("returns expected error", func(t *testing.T) {
		wanted := verr
		got := err
		if !errors.Is(got, wanted) {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}

		t.Run("wrapped in ValidationError", func(t *testing.T) {
			wanted := ValidationError{}
			got := err
			if !errors.As(got, &wanted) {
				t.Errorf("\nwanted %T\ngot    %T", wanted, got)
			}
		})
	})
}
