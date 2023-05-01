package mediator

import (
	"context"
)

// validate calls the supplied Validator for the context specified request.
// Any error returned by the validator is returned; if an error is not a
// ValidationError then it is wrapped in one.
func validate[TRequest any](v Validator[TRequest], ctx context.Context, rq TRequest) error {
	if err := v.Validate(ctx, rq); err != nil {
		if _, ok := err.(ValidationError); ok {
			return err
		}
		return ValidationError{E: err}
	}
	return nil
}
