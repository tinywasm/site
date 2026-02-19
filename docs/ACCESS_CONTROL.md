# Access Control

> **Status:** Current — February 2026

`tinywasm/site` supports two mutually exclusive access control modes, both configured via
the `site` package before calling `RegisterHandlers`.

---

## Mode 1 — Standalone (`SetUserRoles`)

The handler declares which role codes are allowed per action. site extracts the current
user's roles from the request and crudp matches them.

```go
// Configure role extraction from request (call before RegisterHandlers)
site.SetUserRoles(func(data ...any) []byte {
    for _, d := range data {
        if req, ok := d.(*http.Request); ok {
            userID := req.Header.Get("X-User-ID")
            // Simple lookup — could be JWT, session, etc.
            return getUserRoleCodes(userID) // e.g. []byte{'a', 'e'}
        }
    }
    return nil
})

site.RegisterHandlers(modules.Init()...)
```

Each handler declares its own access rules:

```go
func (h *InvoiceHandler) AllowedRoles(action byte) []byte {
    switch action {
    case 'r': return []byte{'*'}            // public read
    case 'c', 'u', 'd': return []byte{'a', 'e'} // admin + editor only
    }
    return nil
}
```

**When to use:** Simple applications without a permission database. Roles are hardcoded
in handler declarations.

---

## Mode 2 — rbac (`SetAccessCheck`)

An external access check function is injected into crudp via site. The handler does **not**
need to implement `AllowedRoles()`. The closure makes the permission decision.

```go
// Wire rbac into site (call before RegisterHandlers)
site.SetAccessCheck(func(resource string, action byte, data ...any) bool {
    for _, d := range data {
        if req, ok := d.(*http.Request); ok {
            userID := req.Header.Get("X-User-ID")
            if userID == "" {
                return false // unauthenticated
            }
            ok, _ := rbac.HasPermission(userID, resource, action)
            return ok
        }
    }
    return false
})

site.RegisterHandlers(modules.Init()...)
```

`tinywasm/rbac` is NOT imported by `tinywasm/site`. The closure belongs to the application
layer. See [ARCHITECTURE.md](ARCHITECTURE.md) for the full startup sequence with rbac.

**When to use:** Applications where permissions are stored in a database and managed at
runtime (create/revoke roles, assign/revoke permissions). Required when roles need to be
configurable without code changes.

---

## Comparison

| | `SetUserRoles` (standalone) | `SetAccessCheck` (rbac) |
|-|--------------------------|------------------------|
| Handler needs `AllowedRoles()` | Yes | No |
| Permissions stored in DB | No | Yes (`rbac_permissions`) |
| Runtime permission changes | No | Yes (rbac API) |
| Extra dependencies | None | `github.com/tinywasm/rbac` (app layer) |
| Complexity | Low | Medium |
| Best for | Simple apps, prototypes | Production multi-role apps |

> **These modes are mutually exclusive.** Calling both `SetUserRoles` and `SetAccessCheck`
> is undefined behavior. Use one or the other.

---

## SSR vs SPA Rendering Decision

The SSR/SPA split is decided per module at startup, based on `AllowedRoles('r')`:

| `AllowedRoles('r')` result | Rendering mode | Effect |
|----------------------------|----------------|--------|
| `[]byte{'*'}` (wildcard) | **SSR** | Full HTML rendered server-side. Indexed by search engines. |
| Any specific roles | **SPA** | Server renders placeholder. WASM client authenticates and renders. |

> In rbac mode (`SetAccessCheck`), if the handler does not implement `AllowedRoles()`,
> site defaults to SPA rendering for that module.

---

## Dev Mode

When `site.SetDevMode(true)` is active, all access checks are bypassed. Never use in production.

```go
site.SetDevMode(true) // all handlers accessible regardless of roles
```
