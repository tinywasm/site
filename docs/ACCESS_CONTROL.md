# Access Control & SSR vs SPA

`tinywasm/site` uses the `AccessLevel` interface from `tinywasm/crudp` to determine both security and the rendering strategy.

## Interface Integration

Components can implement `AllowedRoles(action byte) []byte`.

```go
type AccessLevel interface {
    AllowedRoles(action byte) []byte
}
```

## Rendering Decision Logic

The rendering mode (SSR or SPA) is decided per route based on the allowed roles for the 'read' (`'r'`) action:

| Allowed Roles | Rendering Mode | Reason |
|---------------|----------------|--------|
| Contains `'*'` | **SSR** | Public content should be fast and indexable by search engines. |
| Specific Roles | **SPA / WASM** | Private content requires authentication checks typically handled by the WASM client. |

### SSR (Public)
If a component is public, `site` renders the full HTML on the server. The WASM client can still hydrate it to add interactivity.

### SPA/WASM (Private)
If a component is restricted, the server may render a "loading" or "unauthorized" placeholder (or redirect), and the WASM client handles the authentication and subsequent rendering.

## Mixed Components

A common pattern is a public page (SSR) with private elements.

1. **Server**: Renders the public part of the component. Private parts are rendered as empty containers (slots) with specific IDs.
2. **Client**: The WASM app detects the user session and "mounts" the private components into the reserved slots.

---
**Status**: No Implemented
