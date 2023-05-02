package mediator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMockCommand(t *testing.T) {
	if len(commands) > 0 {
		t.Fatal("invalid test: one or more commands are already registered")
	}

	mock := MockCommand[string, NoResultType]()
	defer mock.Unregister()

	t.Run("registers the mock", func(t *testing.T) {
		wanted := 1
		got := len(commands)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	// ACT
	_, err := Execute(context.Background(), "test", NoResult)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("captures number of calls received", func(t *testing.T) {
		wanted := true
		got := mock.NumRequests() == 1
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("captures copy of requests received", func(t *testing.T) {
		wanted := []string{"test"}
		got := mock.Requests()

		if reflect.ValueOf(got).UnsafePointer() == reflect.ValueOf(mock.requests).UnsafePointer() {
			t.Error("got same slice")
		}

		if !reflect.DeepEqual(wanted, got) {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("when the command is called", func(t *testing.T) {
		wantedWc := true
		wantedWnc := false
		gotWc := mock.WasCalled()
		gotWnc := mock.WasNotCalled()

		if wantedWc != gotWc || wantedWnc != gotWnc {
			t.Errorf("called / not called: wanted %v / %v, got %v / %v", wantedWc, wantedWnc, gotWc, gotWnc)
		}
	})

	t.Run("when the command is not called", func(t *testing.T) {
		mock := MockCommand[int, int]()
		defer mock.Unregister()

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
	herr := errors.New("command error")
	mock := MockCommandError[string, NoResultType](herr)
	defer mock.Unregister()

	// ACT
	_, err := Execute(ctx, "test", NoResult)

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
	defer mock.Unregister()

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

func TestMockCommandValidationError(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	verr := errors.New("validation error")
	mock := MockCommandValidationError[string, NoResultType](verr)
	defer mock.Unregister()

	// ACT
	_, err := Execute(ctx, "test", NoResult)

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

// registermocktestcmd is a mock command for testing the RegisterMock function.
type registermocktestcmd struct {
	cfgerr error
}

func (cmd *registermocktestcmd) CheckConfiguration(context.Context) error { return cmd.cfgerr }
func (*registermocktestcmd) Execute(ctx context.Context, request int) (NoResultType, error) {
	return nil, nil
}

func TestRegisterMock(t *testing.T) {
	// ARRANGE
	ctx := context.Background()

	t.Run("calls the registration function", func(t *testing.T) {
		// ARRANGE
		ofn := register
		defer func() { register = ofn }()

		registerIsCalled := false
		register = func(context.Context, any, any) (func(), error) { registerIsCalled = true; return nil, nil }

		// ACT
		RegisterMockCommand[int, NoResultType](ctx, &registermocktestcmd{})

		// ASSERT
		wanted := true
		got := registerIsCalled
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("registering a custom mock", func(t *testing.T) {
		// ARRANGE
		cmd := &registermocktestcmd{}

		// ACT
		fn := RegisterMockCommand[int, NoResultType](ctx, cmd)

		// ASSERT
		t.Run("adds command to registry", func(t *testing.T) {
			wanted := true
			got := commands[reflect.TypeOf(1)] == cmd
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})

		t.Run("returns non-nil func", func(t *testing.T) {
			wanted := false
			got := fn == nil
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}

			t.Run("which unregisters the custom mock", func(t *testing.T) {
				// ACT
				fn()

				// ASSERT
				wanted := true
				got := commands[reflect.TypeOf(1)] == nil
				if wanted != got {
					t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
				}
			})
		})
	})

	t.Run("panics when registering custom mock with configuration failure", func(t *testing.T) {
		// ARRANGE
		cmd := &registermocktestcmd{cfgerr: errors.New("configuration error")}

		defer func() { // panic tests must be deferred
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		// ACT
		RegisterMockCommand[int, NoResultType](ctx, cmd)
	})
}
