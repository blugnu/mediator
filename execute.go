package mediator

import (
	"context"
	"reflect"
)

// Execute sends the specified request to the registered command for the
// request type. The result parameter is a type-hint, providing a pointer to
// a value of the expected result type.  The result itself is returned by
// the function along with any error.  The type-hint result pointer is
// otherwise ignored.
//
// The result type-hint enables the GoLang compiler to infer both the request
// and result type for the generic function.
//
// The type-hint pointer is not de-referenced so may be nil; alternatively
// the `new()` function may be used to create a pointer to the result type:
//
//	// call the Foo.Request command and capture the result in foo
//	foo, err := mediator.Execute(ctx, Foo.Request{}, new(Foo.Result))
//
// In the event of an error, the result will be the zero-value of the result
// type, otherwise it will be the value returned by the command.
//
// A special case is when a command does not return any result.
//
// Such a command would typically be registered with a result type of
// 'mediator.NoResultType' and called with a type-hint parameter of
// `mediator.NoResult`.  The value returned by the Execute function is
// discarded in this case:
//
//	// call the Bar.Request command which returns only an error
//	_, err := mediator.Execute(ctx, Bar.Request{}, mediator.NoResult)
//
// If the command implements Validator and the validator returns an error,
// then the command Execute() function is not called and the error returned
// will be a ValidationError wrapping the error.
func Execute[TRequest any, TResult any](ctx context.Context, req TRequest, resultHint *TResult) (TResult, error) {
	// create a zero-value result for use in error conditions
	z := *new(TResult)

	// identify the command registration for the request type
	reg, ok := commands[reflect.TypeOf(req)]
	if !ok {
		return z, &NoCommandForRequestTypeError{req}
	}

	// check that registered command returns the result type expected by the caller
	cmd, ok := reg.(CommandHandler[TRequest, TResult])
	if !ok {
		return z, &ResultTypeError{command: reg, result: z}
	}

	// call the Validator, if implemented
	if validator, ok := reg.(Validator[TRequest]); ok {
		err := validate(validator, ctx, req)
		if err != nil {
			return z, err
		}
	}

	// call the command and return the result
	return cmd.Execute(ctx, req)
}
