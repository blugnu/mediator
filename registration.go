package mediator

import (
	"context"
	"reflect"
)

var commands = map[reflect.Type]any{}

// register provides the function used to register a command.  This is called
// by RegisterCommand and the mock command factories.
//
// The function is a variable to facilitate module unit tests.
var register = func(ctx context.Context, rq any, cmd any) (func(), error) {
	rqt := reflect.TypeOf(rq)

	if cmd, exists := commands[rqt]; exists {
		return nil, CommandAlreadyRegisteredError{command: cmd, request: rq}
	}

	// call the ConfigurationChecker, if implemented
	if cfg, ok := cmd.(ConfigurationChecker); ok {
		if err := cfg.CheckConfiguration(ctx); err != nil {
			return nil, err
		}
	}

	commands[rqt] = cmd

	return func() { delete(commands, rqt) }, nil
}

// RegisterCommand[TRequest, TResult] registers a command returning a specific
// result type for the specified request type.
//
// If a command is already registered for the request type the function will
// return a CommandAlreadyRegisteredError.
//
// If the command being registered implements the ConfigurationChecker interface,
// this is called and any error returned.
//
// If the command does not implementation ConfigurationChecker or the configuration
// check returns no error, then the command is registered.
func RegisterCommand[TRequest any, TResult any](ctx context.Context, cmd CommandHandler[TRequest, TResult]) error {
	_, err := register(ctx, *new(TRequest), cmd)
	return err
}
