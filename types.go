package mediator

// NoResult is a convenience type that can be used as the TResult of a command
// when the command does not return a result.
//
// For example:
//
//	// to register a command handler
//	mediator.RegisterCommand[MyRequestType, mediator.NoResult](MyCommandHandler{})
//
//	// calling the command
//	_, err := mediator.Execute(ctx, MyRequestType{}, new(mediator.NoResult))
type NoResult int
