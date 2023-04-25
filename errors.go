package mediator

import (
	"fmt"
)

// NoHandlerError is returned by Execute if there is no handler
// registered for the request and result type involved.
type NoHandlerError struct {
	request any
	result  any
}

func (e *NoHandlerError) Error() string {
	return fmt.Sprintf("no handler for '%T' request returning '%T'", e.request, e.result)
}

func (e *NoHandlerError) Is(target error) bool {
	other, ok := target.(*NoHandlerError)
	return ok && other.request == e.request && other.result == e.result
}

// ResultTypeError is returned by Perform if the registered
// handler for the specified request type does not return then
// specified result type.
type ResultTypeError struct {
	handler interface{}
	request interface{}
	result  interface{}
}

func (e ResultTypeError) Error() string {
	return fmt.Sprintf("handler for %T (%T) does not return %T", e.request, e.handler, e.result)
}

// ValidationError is returned by a handler when it is unable to
// process a request due to the request itself being invalid.  The
// ValidationError wraps a specific error that identifies the
// problem with the request.
type ConfigurationError struct {
	handler any
	E       error
}

func (e ConfigurationError) Error() string {
	return fmt.Sprintf("%T configuration error: %v", e.handler, e.E)
}

func (e ConfigurationError) Unwrap() error {
	return e.E
}

// ValidationError is returned by a handler when it is unable to
// process a request due to the request itself being invalid.  The
// ValidationError wraps a specific error that identifies the
// problem with the request.
type ValidationError struct {
	E error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.E)
}

func (e ValidationError) Unwrap() error {
	return e.E
}
