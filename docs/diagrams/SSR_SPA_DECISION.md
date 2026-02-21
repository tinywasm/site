# SSR vs SPA Rendering Decision

> **Status:** Current — February 2026

`RegisterHandlers` inspects `AllowedRoles('r')` on each handler to determine the rendering
mode. This decision is made once at startup — not per request.

```mermaid
flowchart TD
    A[RegisterHandlers called] --> B[For each handler]
    B --> C{Implements\nAllowedRoles?}
    C -- No --> SPA[SPA mode\nWASM renders after auth]
    C -- Yes --> D[AllowedRoles 'r']
    D --> E{Result}
    E -- "[]byte{'*'}" --> SSR[SSR mode\nHTML injected at startup]
    E -- specific roles\ne.g. 'a','e' --> SPA
    E -- nil --> SPA

    SSR --> F[ssrBuild: RenderPage → assetmin\nHTML embedded in JS bundle\nVisible to search engines]
    SPA --> G[Placeholder rendered\nWASM client fetches after auth\nContent hidden until authenticated]

    B --> B
```

## Effect per mode

| Condition | Mode | First paint | SEO | Auth required |
|-----------|------|-------------|-----|---------------|
| `AllowedRoles('r') == []byte{'*'}` | **SSR** | Instant (HTML in bundle) | Indexed | No |
| `AllowedRoles('r') == []byte{'a','e'}` | **SPA** | After WASM load + auth | Not indexed | Yes |
| `AllowedRoles('r') == nil` | **SPA** | After WASM load + auth | Not indexed | Yes |
| Handler has no `AllowedRoles` method | **SPA** | After WASM load + auth | Not indexed | Yes |

## Code reference

`isPublicReadable(h any) bool` in `register_ssr.go` — reads `AllowedRoles('r')` and
checks for `[]byte{'*'}` wildcard.

## Tests

| Test | Branch covered |
|------|---------------|
| `TestSSRDecision_Wildcard` | `AllowedRoles('r')==[]byte{'*'}` → SSR |
| `TestSSRDecision_SpecificRoles` | `AllowedRoles('r')==[]byte{'a'}` → SPA |
| `TestSSRDecision_Nil` | `AllowedRoles('r')==nil` → SPA |
| `TestSSRDecision_NoInterface` | handler without AllowedRoles → SPA |
