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

## Module Navigation

While `crudp` handles data routes, `tinywasm/site` provides top-level module navigation using URL hashes.

```go
// Navigate to the 'users' module
site.Navigate("app", "users")
```

This updates `window.location.hash` to `#users`, unmounts the previous module, and mounts/hydrates the new one.

---
**Status**: Implemented
