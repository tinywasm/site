# Routing

Routing in `tinywasm/site` is powered by `tinywasm/crudp`. Every component's routes and behavior are defined by the interfaces it implements and its `HandlerName()`.

## Interface-Driven Routes

A component becomes a router by implementing one or more of the `crudp` interfaces. The mapping between methods and HTTP verbs/paths is automatic:

| Interface | Method | HTTP Method | Route Context |
|-----------|--------|-------------|---------------|
| `Creator` | `Create` | `POST` | `/{handler_name}/` |
| `Reader`  | `Read`   | `GET`  | `/{handler_name}/{path...}` |
| `Updater` | `Update` | `PUT`  | `/{handler_name}/{path...}` |
| `Deleter` | `Delete` | `DELETE`| `/{handler_name}/{path...}` |

### Example: User Component

```go
type User struct {}

func (u *User) HandlerName() string { return "users" }

// GET /users/123
func (u *User) Read(data ...any) any {
    // Logic for rendering or returning data
    return u
}
```

## Automatic Registration

When `site.AddModules(components...)` is called, it iterates through all provided components and registers them with the underlying `crudp` router for both Server and WASM.

- **Server-side**: Generates HTTP endpoints via `RegisterRoutes(mux)`.
- **Client-side**: Configures the SPA navigator to handle these routes asynchronously.

## Shared Domain Logic

---
**Status**: Partially Implemented
