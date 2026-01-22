# Modules

Modules in `tinywasm/site` are logical groupings of components. They serve as the entry point for registering related functionality.

## Module Structure

A module should provide an `Add()` method that returns a slice of components (or other handlers).

```go
// modules/users/users.go
type UsersModule struct{}

func (m *UsersModule) Add() []any {
    return []any{
        &UserList{},    // Component 1
        &UserProfile{}, // Component 2
    }
}
```

## Centralized Initialization

All modules are initialized in a central location, typically `modules/init.go`.

```go
// modules/init.go
func Init() []any {
    users := &users.UsersModule{}
    contact := &contact.ContactModule{}
    
    var all []any
    all = append(all, users.Add()...)
    all = append(all, contact.Add()...)
    return all
}
```

## Component-Based Routing

Each component returned by a module's `Add()` method defines its own route by implementing the `NamedHandler` interface from `tinywasm/crudp`.

```go
func (c *UserList) HandlerName() string { return "users" } // Route: /users
func (c *UserProfile) HandlerName() string { return "profile" } // Route: /profile
```

### Typed Routing (Avoid Hardcoded Strings)

To navigate between components without hardcoding strings, use the component's type or instance to resolve its route:

```go
// Conceptual API
url := site.LinkTo(&UserProfile{}) // Returns "/profile"
```

---
**Status**: Partially Implemented

---
**Status**: No Implemented
