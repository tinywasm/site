# Modules

Modules are top-level navigable units in `tinywasm/site`. Each module maps to a hash route (`#handler-name`) and preserves its state across navigations using an internal LRU cache.

## Module Interface (required)

Every module must implement the `site.Module` interface:

```go
type Module interface {
    dom.Component               // RenderHTML(), OnMount()
    HandlerName() string        // route key (e.g. "users" → "#users")
    ModuleTitle() string        // display title
}
```

## Optional Interfaces

Modules can extend their behavior by implementing additional interfaces:

| Interface | Method(s) | When to implement |
|-----------|-----------|-------------------|
| `Parameterized` | `SetParams(params []string)` | When the route carries path params (e.g. `#users/123`) |
| `ModuleLifecycle` | `BeforeNavigateAway() bool`<br>`AfterNavigateTo()` | Navigation hooks (block navigation, refresh on return) |
| `CSSProvider` | `RenderCSS() string` | Per-module CSS injected at startup (!wasm) |
| `JSProvider` | `RenderJS() string` | Per-module JS injected at startup (!wasm) |
| `IconSvgProvider` | `IconSvg() map[string]string` | SVG sprite entry for this module (!wasm) |
| `AccessLevel` | `AllowedRoles(action byte) []byte` | RBAC: declares which roles can access each action |

## Registration

Modules are typically initialized in a dedicated `modules` package and registered at the application start:

```go
hs := modules.Init()  // returns []any of all handlers/modules
site.RegisterHandlers(hs...)
```

## Hash Routing

The `HandlerName()` is used as the routing key. The `site` package normalizes the URL hash:
- `"users"`, `"#users"`, `"#/users"` → module `users`, no params
- `"#users/123"` → module `users`, params `["123"]`

## Path Params Example

When a module implements `Parameterized`, `SetParams` is called automatically upon navigation:

```go
func (u *Users) SetParams(params []string) {
    if len(params) > 0 {
        u.selectedID = params[0]
    }
}
```

Navigating with params:
```go
site.Navigate("app", "users/123")  // Calls SetParams(["123"]) on Users module
```

## Lifecycle Hooks

Use `ModuleLifecycle` to control navigation or refresh state:

```go
func (u *Users) BeforeNavigateAway() bool {
    return !u.hasUnsavedChanges  // returns false to block navigation
}

func (u *Users) AfterNavigateTo() {
    u.Refresh()  // called when user navigates back from cache or initial load
}
```

## Module Manager

- **LRU Cache**: The site maintains a cache (default size: 3) to preserve Go struct state. This enables instant "back" navigation.
- **`site.SetCacheSize(n)`**: Configures the number of cached modules.
- **`site.Mount(parentID string)` (wasm)**: Hydrates the initial module from the URL hash and starts the main loop. Blocks forever.
