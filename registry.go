package mediator

import (
	"fmt"
	"reflect"
)

var commandHandlers = map[reflect.Type]any{}

// RegisterCommand registers a handler for the specified request type
// returning the specified result type.
//
// If a handler is already registered for the request type the
// function will panic, otherwise the handler is registered.
func RegisterCommand[TRequest any, TResult any](handler CommandHandler[TRequest, TResult]) func() {
	rq := *new(TRequest)
	rs := *new(TResult)

	rqt := reflect.TypeOf(rq)

	_, exists := commandHandlers[rqt]
	if exists {
		panic(fmt.Sprintf("command already registered for %T request returning %T", rq, rs))
	}

	commandHandlers[rqt] = handler

	return func() { delete(commandHandlers, rqt) }
}
