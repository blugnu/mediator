package mediator

import (
	"fmt"
	"reflect"
)

// NoCommandForRequestTypeError is returned by Execute if there is no command
// registered for the request and result type involved.
type CommandAlreadyRegisteredError struct {
	command any
	request any
}

func (e CommandAlreadyRegisteredError) Error() string {
	return fmt.Sprintf("%T already registered for requests of type: %T", e.command, e.request)
}

func (e CommandAlreadyRegisteredError) Is(target error) bool {
	if other, ok := target.(CommandAlreadyRegisteredError); ok {
		return ok && reflect.TypeOf(other.request) == reflect.TypeOf(e.request)
	}
	if other, ok := target.(*CommandAlreadyRegisteredError); ok {
		return ok && reflect.TypeOf(other.request) == reflect.TypeOf(e.request)
	}
	return false
}

// NoCommandForRequestTypeError is returned by Execute if there is no command
// registered for the request and result type involved.
type NoCommandForRequestTypeError struct {
	request any
}

func (e NoCommandForRequestTypeError) Error() string {
	return fmt.Sprintf("no command registered for requests of type: %T", e.request)
}

func (e NoCommandForRequestTypeError) Is(target error) bool {
	if other, ok := target.(NoCommandForRequestTypeError); ok {
		return ok && reflect.TypeOf(other.request) == reflect.TypeOf(e.request)
	}
	if other, ok := target.(*NoCommandForRequestTypeError); ok {
		return ok && reflect.TypeOf(other.request) == reflect.TypeOf(e.request)
	}
	return false
}

// ResultTypeError is returned if the command registered for the
// specified request type does not return the result type expected
// by the caller.
type ResultTypeError struct {
	command any
	result  any
}

func (e ResultTypeError) Error() string {
	return fmt.Sprintf("%T does not return %T", e.command, e.result)
}

func (e ResultTypeError) Is(target error) bool {
	if other, ok := target.(ResultTypeError); ok {
		return ok && reflect.TypeOf(other.result) == reflect.TypeOf(e.result)
	}
	if other, ok := target.(*ResultTypeError); ok {
		return ok && reflect.TypeOf(other.result) == reflect.TypeOf(e.result)
	}
	return false
}

// ValidationError is returned by a command when it is unable to
// process a request due to the request itself being invalid.  The
// ValidationError wraps a specific error that identifies the
// problem with the request.
//
//	"request validation error: <specific error>"
type ValidationError struct {
	E error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("request validation error: %v", e.E)
}

func (e ValidationError) Unwrap() error {
	return e.E
}
