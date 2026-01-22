# Components

Components are the fundamental UI and logic units. They integrate `tinywasm/dom` for the view and `tinywasm/crudp` for the behavior and routing.

## Interfaces

A component's capabilities are determined by the interfaces it implements:

### Behavior & Routing (via `crudp`)
- `NamedHandler`: `HandlerName() string` (Required for routing).
- `Reader`, `Creator`, `Updater`, `Deleter`: Map to GET, POST, PUT, DELETE.
- `AccessLevel`: `AllowedRoles(action byte) []byte` (Security & SSR mode).

### UI & Assets (via `site`)
- `RenderHTML() string`: The component's HTML structure.
- `RenderCSS() string`: Component-specific CSS.
- `IconSvg() []map[string]string`: Multi-icon declarations for the global SVG sprite.

### Multi-Icon Support

The `IconSvg()` method returns a slice of maps, allowing a component to declare multiple icons that will be bundled into the global SVG sprite.

```go
func (c *MyComponent) IconSvg() []map[string]string {
    return []map[string]string{
        {"id": "edit-icon", "svg": "<svg>...</svg>"},
        {"id": "save-icon", "svg": "<svg>...</svg>"},
    }
}
```

## Internal logic

Each component can have its own internal state and event handling via `OnMount()`, as defined in [tinywasm/dom components](github.com/tinywasm/dom/blob/main/docs/COMPONENTS.md).

---
**Status**: No Implemented
