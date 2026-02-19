# Components

Components are the fundamental UI and logic units. They integrate `tinywasm/dom` for the view and `tinywasm/crudp` for the behavior and routing.

## Interfaces

A component's capabilities are determined by the interfaces it implements. These interfaces are checked via type assertion at registration time.

### Behavior & Routing (via `crudp`)

- `NamedHandler`: `HandlerName() string` (Required for routing).
- `Reader`, `Creator`, `Updater`, `Deleter`: Direct CRUD mapping to HTTP verbs.
- `AccessLevel`: `AllowedRoles(action byte) []byte` (RBAC & SSR/SPA decision).
- `DataValidator`: `ValidateData(action byte, data ...any) error`.

### UI & Assets (via `site`)

- `dom.Component`: `RenderHTML() string` and `OnMount()`.
- `CSSProvider`: `RenderCSS() string`.
- `JSProvider`: `RenderJS() string`.
- `IconSvgProvider`: `IconSvg() map[string]string`.

## Asset Providers

Components can provide their own styles, logic, and icons to be bundled into the global assets.

### CSS & JS
Implement `CSSProvider` or `JSProvider` to inject raw CSS/JS into the global bundles. This is handled by `tinywasm/site` at registration time (backend only).

### Icons
The `IconSvg()` method returns a single map containing the icon ID and its SVG source. This icon will be bundled into the global SVG sprite.

```go
func (c *MyComponent) IconSvg() map[string]string {
    return map[string]string{
        "id":  "edit-icon",
        "svg": "<svg>...</svg>",
    }
}
```

## Internal logic

Each component can have its own internal state and event handling via `OnMount()`, as defined in [tinywasm/dom components](github.com/tinywasm/dom/blob/main/docs/COMPONENTS.md).
