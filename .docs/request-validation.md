# Request Validation

Commands typically perform some validation of any request, to ensure that the request is valid.

Request validation may be performed either in-line in the command Execute function itself, or separated into a `Validate()` function.

A `ValidationError` type is provided, to enable callers to differentiate between an incorrectly specified request and an error in the execution of the command to satisfy a particular request. 

1. Implementation Options
2. Implementation Examples
3. The `mediator.ValidationError` Type
4. Separation of Concerns

<br/>

## Implementation Options

Validation of a command request may be implemented in one of two ways:

1. As part of the `Execute` function (`CommandHandler` interface)
2. In a separate `Validate` function (`Validator` interface)


## Implementation Examples

### As part of the `Execute` function:

```golang
func (h *Handler) Execute(ctx context.Context, rq Request) (*Result, error) {
    if !h.Exists(ctx, rq.Id) {
        return nil, mediator.ValidationError{ErrDoesNotExist}
    }

    // command execution continues...
}
```

### As a `Validate` function:

> _Any error returned from the `Validate` function is **automatically wrapped** in a `mediator.ValidationError`; `Validate` should simply return the detailed error itself._

```golang
func (h *Handler) Validate(ctx context.Context, rq Request) error {
    if !h.Exists(ctx, rq.Id) {
        return ErrDoesNotExist
    }
    return nil
}
```

## The `mediator.ValidationError` type

A command should ensure that errors relating to request validation are returned as a `mediator.ValidationError` wrapping a specific error identifying the validation failure.

This enables callers to differentiate between invalid requests (_a mistake made by the caller_) and errors in the execution of the command itself.

For example, this enables HTTP endpoints calling commands via the mediator to determine when a `400 bad request` is a more appropriate response than `500 internal server error`.

<br/>

# Separation of Concerns

Consider a command request that is triggered in response to some HTTP request received at a REST Api endpoint.

The HTTP endpoint must first ensure that _the HTTP request_ itself satisfies the endpoint contract.  This might include:

- ensuring that the expected/supported HTTP method has been used
- ensuring that required parameters in the resource url or query string are present and have valid/expected types or values
- etc.

If the incoming HTTP request does not satisfy the _endpoint_ contract then the endpoint will reject the request without ever calling any command.

Otherwise the endpoint will intialise a command request and call the command via the mediator.

The command must now ensure that this request is valid.

> _Command request validation should **not** be concerned with, or try to direct, how any caller should respond to any particular form of invalid request.  Appropriate, specific errors should be returned to enable the caller to differentiate between different scenarios and respond accordingly._


## The Separate Concerns

- an HTTP Api is concerned with ensuring that a valid HTTP request is received before being passed to the command.

- _an HTTP Api is **not** concerned with whether the received request makes sense to and is actionable by the command_.

- a command is concerned with ensuring that a particular request is valid and able to be executed.

- _a command is **not** concerned with how any caller will respond to any particular request validation error; it must provide appropriate error responses to enable a caller to make any appropriate distinctions_.


## Other Scenarios

Similar distinctions apply when there is some other stimulus for a command request.

For example, if the stimulus is a Kafka event then instead of HTTP request validation there will be some message handler that validates a Kafka message payload.
