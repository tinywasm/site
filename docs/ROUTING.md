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

## Module Navigation (WASM)

While `crudp` handles data routes, `tinywasm/site` provides top-level module navigation using URL hashes.

```go
// site.Navigate(parentID, hash string)
site.Navigate("app", "users/123")
```

The `hash` parameter is normalized by `parseRoute()`. All these formats are accepted:

| Hash String | Resolved Module | Params |
|-------------|-----------------|--------|
| `"users"` | `users` | `[]` |
| `"#users"` | `users` | `[]` |
| `"#/users"` | `users` | `[]` |
| `"#users/123"` | `users` | `["123"]` |

Navigation updates `window.location.hash`, unmounts the current module, and mounts/hydrates the new one from the LRU cache if available.
