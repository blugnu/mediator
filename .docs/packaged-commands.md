# Packaged Commands

1. encourages separation of concerns between your command handlers
2. encourages adherence to the single responsibility principle
3. enables you to use types called `Request`, and `Handler` for **_all_** your commands (and `Result` where desired)
4. result type names are scoped to the command package

Package qualification of the request, handler and result types result in implementation, registration and execution code that reads naturally.

## Implementation

```golang
package getFoo

// Request a Foo by Id; returns *getFoo.Result or nil
type Request struct {
    Id string
}

type Result struct {
    Id   string `json:"id"`
    Name string `json:"name"`
}

// Handler implements the getFoo command
type Handler struct {
    repository.GetFooByIdMethod   // injectable repository method dependency (implements GetFooById)
}

func (h *Handler) Execute(ctx context.Context, rq Request) (*Result, error) {
    // retrieves the foo with the requested id from the repository
    // using the anonymous GetFooByIdMethod injected into the handler
    foo, err := h.GetFooById(ctx, rq.Id)
    if err != nil {
        return nil, err
    }

    // marshal the repository.Foo model to a getFoo.Result
    result := &Result{
        Id: foo.Id,
        Name: foo.NamePrefix + foo.Name,
    }
    return result, nil
}
```
