package mediator

import (
	"context"
	"reflect"
)

// CommandHandler[TRequest, TResult] is the interface that must be implemented
// by a command handler.
type CommandHandler[TRequest any, TResult any] interface {
	Execute(context.Context, TRequest) (TResult, error)
}

// Execute sends the specified request and context to the registered command
// handler for the request type, providing a pointer to a variable to receive
// the result. The result is also be returned by the Execute function, together
// with any error returned by the handler.
//
// If the handler implements Validator and the validator returns an error,
// then the handler is not called and the error returned will be a
// ValidationError, wrapping the error returned by the validator.
//
// The pointer to the result variable allows GoLang to infer both the request
// and result type for the generic function, so these do not need to be
// specified explicitly.
//
// Because the result value is also returned from the function you can use
// 'new(<RequestType>)' to pass a new pointer, assigning the function result
// to a new variable which will be of the appropriate result type:
//
//	// call the FooRequest handler and place result in new variable
//	foo, err := mediator.Execute(ctx, FooRequest{}, new(FooResult))
//
// If you have a variable already declared that you wish to receive the
// result, you can pass a pointer to that variable instead and discard the
// returned result.  This can be useful to create a conditionally scoped 'error':
//
//	// call the FooRequest handler and place result in existing variable
//	var foo = FooResult
//	if _, err := mediator.Execute(ctx, FooRequest{}, &foo); err != nil {
//	  handle error
//	}
//
// A special case is when the command does not return any result.  Such a
// command handler would typically be registered with a result type of
// 'mediator.NoResult'.  A pointer of this type must still be supplied; the
// value returned by the Execute function should be discarded:
//
//	// call the BarRequest handler which returns only an error
//	_, err := mediator.Execute(ctx, BarRequest{}, new(mediator.NoResult))
func Execute[TRequest any, TResult any](ctx context.Context, rq TRequest, rs *TResult) (result TResult, err error) {
	rqt := reflect.TypeOf(rq)

	reg, ok := commandHandlers[rqt]
	if !ok {
		return *rs, &NoHandlerError{rq, *rs}
	}

	handler, ok := reg.(CommandHandler[TRequest, TResult])
	if !ok {
		return *rs, &ResultTypeError{handler: handler, request: rq, result: *rs}
	}

	// If the handler implements ConfigurationChecker call that first
	// and return any error
	if validator, ok := reg.(ConfigurationChecker[TRequest]); ok {
		err := checkConfiguration(validator, ctx)
		if err != nil {
			return *rs, err
		}
	}

	// If the handler implements validator, call that first
	// and return any error
	if validator, ok := reg.(Validator[TRequest]); ok {
		err := validate(validator, ctx, rq)
		if err != nil {
			return *rs, err
		}
	}

	*rs, err = handler.Execute(ctx, rq)
	return *rs, err
}
