package mediator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// registrationtestcmd is a command used for testing the registration function.
type registrationtestcmd struct {
	cfgerr error
}

func (cmd registrationtestcmd) CheckConfiguration(context.Context) error       { return cmd.cfgerr }
func (registrationtestcmd) Execute(context.Context, int) (NoResultType, error) { return nil, nil }

func TestRegistrationFunction(t *testing.T) {
	// ARRANGE
	ctx := context.Background()

	t.Run("when command already registered for request type", func(t *testing.T) {
		// ARRANGE
		commands[reflect.TypeOf(1)] = registrationtestcmd{}
		defer func() { delete(commands, reflect.TypeOf(1)) }()

		// ACT
		fn, err := register(ctx, *new(int), registrationtestcmd{})

		// ASSERT
		t.Run("returns nil func", func(t *testing.T) {
			wanted := true
			got := fn == nil
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})

		t.Run("returns error", func(t *testing.T) {
			wanted := CommandAlreadyRegisteredError{command: registrationtestcmd{}, request: *new(int)}
			got := err
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	})

	t.Run("returns any ConfigurationChecker error", func(t *testing.T) {
		// ARRANGE
		cfgerr := errors.New("configuration error")
		cmd := registrationtestcmd{cfgerr}

		// ACT
		fn, err := register(ctx, *new(int), cmd)

		// ASSERT
		t.Run("returns nil func", func(t *testing.T) {
			wanted := true
			got := fn == nil
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})

		t.Run("returns error", func(t *testing.T) {
			wanted := cfgerr
			got := err
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	})

	t.Run("registering a command", func(t *testing.T) {
		// ARRANGE
		cmd := registrationtestcmd{}

		// ACT
		fn, err := register(ctx, *new(int), cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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

			t.Run("which unregisters the command", func(t *testing.T) {
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
}

func TestRegisterCommand(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	cmd := registrationtestcmd{}

	t.Run("calls register function", func(t *testing.T) {
		// ARRANGE
		ofn := register
		defer func() { register = ofn }()

		registerIsCalled := false
		register = func(context.Context, any, any) (func(), error) { registerIsCalled = true; return nil, nil }

		// ACT
		err := RegisterCommand[int, NoResultType](ctx, cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// ASSERT
		wanted := true
		got := registerIsCalled
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("returns register function error", func(t *testing.T) {
		// ARRANGE
		ofn := register
		defer func() { register = ofn }()

		regerr := errors.New("registration error")
		register = func(context.Context, any, any) (func(), error) { return nil, regerr }

		// ACT
		err := RegisterCommand[int, NoResultType](ctx, cmd)

		// ASSERT
		wanted := regerr
		got := err
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}
