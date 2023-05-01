package mediator

import "context"

// NoResultType may be used as the TResult of a command when the command does not
// return a result.  The NoResult value may then be used in Execute calls.
//
// For example:
//
//	// to register a command with no result
//	mediator.RegisterCommand[MyCommand.Request, mediator.NoResultType](MyCommand.Handler{})
//
//	// calling the command
//	_, err := mediator.Execute(ctx, MyCommand.Request{}, mediator.NoResult)
type NoResultType *int

// NoResult may be used in Execute calls to commands that do not return a result:
//
//	// calling the command
//	_, err := mediator.Execute(ctx, MyCommand.Request{}, mediator.NoResult)
//
// Since NoResultType is itself a pointer, when implementing a command that does not
// return a result simply return a nil result:
//
//	// implementing a command returning NoResultType
//	func (c Command) Execute(ctx context.Context, req Request) (mediator.NoResultType, error) {
//	  if err := c.doSomething(); err != nil {
//	    return nil, err
//	  }
//	  return nil, nil
//	}
var NoResult = new(NoResultType)

// CommandHandler[TRequest, TResult] is the one interface that MUST be
// implemented by a command.
//
// It provides the Execute function that is called by the mediator.
type CommandHandler[TRequest any, TResult any] interface {
	Execute(context.Context, TRequest) (TResult, error)
}

// ConfigurationChecker is an optional interface that may be implemented
// by a command to separate any configuration checks from validation
// of any specific request or the execution of the command.
//
// If implemented, CheckConfiguration is called when registering a command.
// if an error is returned, command registration fails and the error is
// returned from RegisterCommand.
type ConfigurationChecker interface {
	CheckConfiguration(ctx context.Context) error
}

// Validator[TRequest] is an optional interface that may be implemented
// by a command to separate the validation of the request from the
// execution of the command.
//
// If implemented, the mediator will the Validate function before
// passing a request to the Execute function; any error returned by
// Validate is returned (wrapped in a ValidationError).
type Validator[TRequest any] interface {
	Validate(context.Context, TRequest) error
}
