<div align="center" style="margin-bottom:20px">
  <img src=".assets/banner.png" alt="mediator" />
  <div align="center">
    <a href="https://github.com/blugnu/mediator/actions/workflows/qa.yml"><img alt="build-status" src="https://github.com/blugnu/mediator/actions/workflows/qa.yml/badge.svg?branch=master&style=flat-square"/></a>
    <a href="https://goreportcard.com/report/github.com/blugnu/mediator" ><img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/mediator"/></a>
    <a><img alt="go version >= 1.18" src="https://img.shields.io/github/go-mod/go-version/blugnu/mediator?style=flat-square"/></a>
    <a href="https://github.com/blugnu/mediator/blob/master/LICENSE"><img alt="MIT License" src="https://img.shields.io/github/license/blugnu/mediator?color=%234275f5&style=flat-square"/></a>
    <a href="https://coveralls.io/github/blugnu/mediator?branch=master"><img alt="coverage" src="https://img.shields.io/coveralls/github/blugnu/mediator?style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/blugnu/mediator"><img alt="docs" src="https://pkg.go.dev/badge/github.com/blugnu/mediator"/></a>
  </div>
</div>

<br/>

# mediator

A lightweight implementation of [the Mediator Pattern](https://en.wikipedia.org/wiki/Mediator_pattern) for `GoLang`, inspired by [jbogard's MediatR framework for .net](https://github.com/jbogard/MediatR).

#### Project History

This project was previously known as `go-mediator`.  It has been renamed as `mediator` for consistency with the package name and because all `blugnu` projects are golang; the `go-` prefix was just noise.

At the same time, the project was completely re-written; it now shares little more than the original concept with the previous incarnation.

If you previously imported `go-mediator` you should update your imports to the renamed module.

<hr/>

## The Mediator Pattern
[The Mediator](https://en.wikipedia.org/wiki/Mediator_pattern) is a simple [pattern](https://en.wikipedia.org/wiki/Software_design_pattern) that uses a 3rd-party (the mediator) to facilitate communication between two other parties without either requiring knowledge of each other.

It is a powerful pattern for achieving loosely coupled code.

There are many ways to implement the pattern, from simple `func` pointers to sophisticated and complex messaging systems; `blugnu/mediator` sits firmly at the *simple* end of that spectrum!

## Why Use `mediator`
For code that provides a substantial component of domain behaviour, using `mediator` provides a consistent mechanism for de-coupling, implementing, calling and mocking those components.

## When NOT To Use `mediator`
Often when testing you may find yourself needing to use a function variable so that you can inject a fake or spy function in order to test higher-level code.  `mediator` is not designed or intended to replace this or similar techniques.

## What `mediator` Is NOT
- it is **not** a message queue
- it is **not** asynchronous
- it is **not** complicated!

<br/>

# How It Works

### TL;DR:

Your code registers commands to respond to requests of various types.  Commands are then called by passing requests to the mediator; the mediator lookups up the command that handles that request, calls it and returns the result and any error.

### In Detail

`blugnu/mediator` maintains a registry of commands that respond to requests of a specific type.  As well as responding to a specific request type, each registered command identifies the result type that it returns to any caller.  

**There can be only _one_ command registered for any given request type.**

Commands are registered during initialising of your application using `RegisterCommand`, or by establishing mock commands in tests.  Command configuration checks are performed when registering commands.  The `RegisterCommand` function tests for an implementation of the `ConfigurationChecker` interface (`CheckConfiguration()` function) which is called if present.  If configuration checks return an error, this is returned by the `RegisterCommand` function and the command is not registered.

Registered commands are called indirectly via a generic `mediator.Execute[TRequest, TResult]` function: **the mediator**.

The mediator consults the registered commands to identify the command for the request type involved.  If no command is registered then a `NoCommandForRequestTypeError` is returned.

If a command is identified but the caller and the command do not agree on the result type, a `ResultTypeError` is returned.

If the correct result type is expected, the mediator tests for an implementation of the `Validator` interface (`Validate()` function) which is called if present.  Any error returned from the `Validate()` function is wrapped in a `ValidationError` (if necessary) and returned to the caller.

If there is no `Validator` interface, or the request is validated successfully, the request is passed to the command and the result and any error from the command then returned to the caller.

All of this takes place _synchronously_ as direct function calls.  i.e. if the command panics, the stack will contain a complete path of execution from the caller, thru the mediator to the corresponding command function.

<br/>
<hr/>

# Implementing a Command

1. (_Recommended_): Create a Package for Your Command
2. Declare request, result and command types
3. (_Optional_) Implement the `ConfigurationChecker` interface for the command
4. (_Optional_) Implement the `Validator` interface for the command
5. Implement the `CommandHandler` interface for the command

> 1. There are numerous advantages to implementing each command in its own package.  See [Packaged Commands](.docs/packaged-commands.md) for more details.

> 3. Any configuration checks incorporated in the `Execute` function are performed _for every request_; performing these checks in a `CheckConfiguration()` function (implementing the `ConfigurationChecker` interface) these checks are performed _just once_, at the time of registering the command.  See [Command Configuration Checks]((#command-configuration-checks)) for more information.

> 4. Any request validation is recommended to be performed in a `Validate()` function (implementing the `Validator` interface).  See [Request Validation](.docs/request-validation.md) for more information.

6. Register the command, e.g.:

```golang
    err := mediator.RegisterCommand[myCommand.Request, *myCommand.Result](ctx, &myCommand.Handler{})
```

> Once a command has been registered it _cannot be **un**registered_, i.e. it is not possible to dynamically reconfigure registered commands to respond to requests of a given type with different commands at different times.  _This is by design_.  In contrast, **_mock_** commands _can_ (and _must_) be reconfigured during the execution of different tests, and this _is_ possible (see: [Testing With Mediator](#testing)).

<br/>
<hr/>

# Calling a Command Using `mediator`

The `mediator.Execute` function accepts a `Context`, the request to be executed and a pointer to a value of the result type.  The function _returns_ the result value and any error from the command.

> _The result type pointer is not de-referenced by the mediator and does not receive any result._

The pointer is required only as a type-hint for the compiler so that it can infer the types required by the generic `Execute` function.

It is recommended to use `new()` to provide a pointer of the required type

#### `example`
```golang
    rq := myCommand.Request{Id: id}
    rs, err := mediator.Execute(ctx, rq, new(*myCommand.Result))
```

> In the above example, `myCommand` returns a pointer to a `myCommand.Result`; `new()` in this case is used to return _a pointer to a pointer_.

## Commands Returning No Result

For commands that have no result value `mediator` provides a convenience type for use when [implementing and registering commands returning no result](#implementing-no-result), and a variable for use as a type-hint when [calling such a command](#calling-no-result):

```golang
    type NoResultType *int
    var NoResult = new(NoResultType)
```

<a name="implementing-no-result"></a>
A command that specifically has no result value is registered with a result type of `mediator.NoResultType` and, as you would expect,  the `Execute()` function of that command returns `mediator.NoResultType`.

> `NoResultType` is a _pointer_ so that when implementing the `Execute()` function for a command returning `NoResultType` you can return `nil`.

#### `example`
```golang
    // Registering a command returning no result
    err := mediator.RegisterCommand[MyRequestType, mediator.NoResultType](ctx, MyCommandHandler{})

    // Implementing the Execute function of a command returning no result
    func (cmd *Handler) Execute(ctx context.Context, req Request) (mediator.NoResultType, error) {
        if err := SomeOperation(); err != nil {
            return nil, err
        }
        return nil, nil
    }
```

<a name="calling-no-result"></a>
A caller can use either `new(mediator.NoResultType)` or `mediator.NoResult` as the result type-hint for the `Execute` function, discarding the returned result.

#### `example`
```golang
    rq := deleteFoo.Request{Id: id}

    // these two statements are functionally equivalent
    _, err := mediator.Execute(ctx, rq, mediator.NoResult)
    _, err := mediator.Execute(ctx, rq, new(mediator.NoResultType))
```

<br/>
<hr/>

# Command Configuration Checks

Before executing any request, a command will typically check the configuration of the command, e.g. to ensure that any required dependencies have been supplied.  This incurs the overhead of those configuration checks on every request when they typically only need to be performed once.

To perform these checks only once, a command may implement the `ConfigurationChecker` interface:

```golang
type ConfigurationChecker interface {
    CheckConfiguration(context.Context) (err error)
}
```

If implemented, the `CheckConfiguration` function is called when _registering_ the command.  If an error is returned from the function then the command registration fails and the error is returned from the `RegisterCommand` function.

<br/>
<hr/>

# Testing With Mediator <a name="testing"></a>

The loose-coupling that can be achieved with a mediator is particularly useful for unit testing.

When unit testing code that calls some command using mediator you are able to mock responses to the request to test the behaviour of your code under a variety of error or result conditions, without having to modify the code under test.


## Mock commands
You can implement mock commands for your request as needed, or you can use the mock factories provided by `blugnu/mediator`; these should be sufficient for most - if not all - common use cases.

The mocks returned by these factories provide an `Unregister()` method to remove the registration for that command; typically you would defer a call to this `Unregister()` method immediately after initialising the mock

#### `example`
```go
    mock := mediator.MockCommand[myCommand.Request, myCommand.Result]()
    defer mock.Unregister()
```

The example above illustrates the mock factory that initialises a command that mocks a successful call, returning a zero-value result and nil error.

The factory functions are:

```golang
    // Mocks a command returning a zero-value result and nil error
    MockCommand[TRequest, TResult]() *mockcommand[TRequest, TResult]

    // Mocks a command returning a specific result and nil error
    MockCommandResult[TRequest, TResult](result TResult) *mockcommand[TRequest, TResult]

    // Mocks a command returning a specific error
    MockCommandError[TRequest, TResult](error) *mockcommand[TRequest, TResult]

    // Mocks a command returning an error from an implementation
    // of the Validator interface
    MockCommandValidationError[TRequest, TResult](error) *mockcommand[TRequest, TResult]
```

> There is no factory for mocking a command that returns an error from a `ConfigurationChecker` interface; such a command would be impossible to register and so could not be called in any test scenario.

The mock returned by these factories provide methods for determining how many times the mock was called, whether it was called at all, as well as copies of all requests received by the mock over its lifetime.

## Custom Mocks

If the provided mock factories are not sufficient, you can register a custom mock using the `RegisterMockCommand()` function.  This is similar to the `RegisterCommand()` function, registering the specified command to handle requests of a specified type and returning a specified result type.

There are two main differences:

- `RegisterMockCommand()` does **not** return any error; if the supplied mock returns an error from any configuration checks, the mock will not be registered and the function will `panic`.
- `RegisterMockCommand()` returns a function to be used to unregister the mock when no longer required (typically immediately deferred to clean up the registration when the test completes)

#### `example`
```golang
    unreg := RegisterMockCommand[myCommand.Request, NoResultType](ctx, &mockMycommand{})
    defer unreg()
```