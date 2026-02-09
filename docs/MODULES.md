# Modules

Modules in `tinywasm/site` are top-level components that orchestrate the site's functionality. They integrate with the `dom` package for rendering and life-cycle management.

## Module Interface

A module must implement the `site.Module` interface:

```go
type Module interface {
    dom.Component
    HandlerName() string // Unique identifier used for routing
    ModuleTitle() string // Title displayed in the site
}
```

## Registration

Register modules (and any other handlers) in your application's entry point:

```go
func init() {
    site.RegisterHandlers(&MyModule{})
}
```

## Lifecycle & Caching

The `site` package managed the active module and maintains a cache of the last 3 active modules (LRU).

- **`site.Start(parentID)`**: Initializes the site, hydrates the initial module based on the URL hash.
- **`site.Navigate(parentID, name)`**: Switches between modules, handling `Unmount` and `Mount` automatically.

Modules preserve their state (as Go structs) when moved to the cache, allowing for instant "back" navigation without data loss.

---
**Status**: Implemented
