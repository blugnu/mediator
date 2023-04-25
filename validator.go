package mediator

import (
	"context"
)

// Validator[TRequest] is an optional interface that may be implemented
// by a command handler to separate the validation of the request from the
// execution of the handler.
//
// This is most useful when the required validation is relatively complex,
// degrading the signal to noise ratio of the Execute function.
//
// Any error returned from the Validate function will be automatically
// wrapped in a ValidationError when returned from the command by the mediator.
// If validation is performed in-line in the Execute function, any validation
// errors must be explicitly wrapped in a ValidationError.
type Validator[TRequest any] interface {
	Validate(context.Context, TRequest) error
}

// validate calls the supplied Validator for the context and input specified.
// Any error returned by the validator is returned; if an error is not a
// ValidationError then it is wrapped in a ValidationError.
func validate[TRequest any](v Validator[TRequest], ctx context.Context, rq TRequest) error {
	if err := v.Validate(ctx, rq); err != nil {
		if _, ok := err.(ValidationError); ok {
			return err
		}
		return ValidationError{err}
	}
	return nil
}
