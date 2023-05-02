# Packaged Commands

1. encourages separation of concerns between your commands
2. encourages adherence to the single responsibility principle
3. enables you to use types called `Request`, and `Handler` (or `Command` etc) for **_all_** your commands (and `Result` where desired)
4. result type names are scoped to the command package

Package qualification of the request, command and result types result in implementation, registration and execution code that _reads_ naturally.

## Implementation

```golang
package getFoo

// Request a Foo by Id; returns *getFoo.Result or nil
type Request struct {
    Id string
}

type Result struct {
    Id   string `json:"id"`
    Name string `json:"name"`   // The full name of the Foo, comprising any prefix and the internal name
}

// Command implements the getFoo command
type Command struct {
    repository.Repository
}

func (cmd *Command) Execute(ctx context.Context, rq Request) (*Result, error) {
    // retrieves the foo with the requested id from the injected, anonymous
    // Repository, providing a GetFooById function which we can use as if it
    // were provided by the Command itself
    //
    // NOTE: the repository method returns a repository model for the 'Foo' ...
    foo, err := cmd.GetFooById(ctx, rq.Id)
    if err != nil {
        return nil, err
    }

    // ... which we must map to a getFoo.Result
    result := &Result{
        Id: foo.Id,
        Name: foo.NamePrefix + foo.InternalName,
    }
    return result, nil
}
```

## Command or Handler or ??

In the example above, the type that implements the command has been named `Command`.  It would equally have been named `Handler` or anything else.  The only code that references this type is the code which initialises and registers your commands; callers never reference it.

In the mediator implementation and documentation, the terms `handler` and `command` are used more-or-less interchangeably.