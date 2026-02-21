# Register Flow

> **Status:** Current — February 2026

`RegisterHandlers` queues work; `applyRBAC` (called at `Mount`/`Serve`) flushes the queue.
This decoupling allows any call order before `Serve`.

```mermaid
sequenceDiagram
    participant App
    participant site
    participant crudp
    participant rbac
    participant assetmin

    Note over App,assetmin: Config phase (any order before Serve)
    App->>site: RegisterHandlers(handlers...)

    site->>site: for each handler: registerModule(m) if Module
    site->>crudp: cp.RegisterHandlers(handlers...)
    Note over crudp: routes + CRUD wired\nAllowedRoles read for access check

    site->>site: registerRBAC(handlers...)
    Note over site: pendingHandlers = append(pendingHandlers, handlers...)

    site->>site: registerAssets(handlers...) [!wasm]
    Note over site: CSS/JS/SVG extracted into ssrState

    Note over App,assetmin: Serve phase
    App->>site: Serve(":8080")
    site->>site: Mount(mux)
    site->>site: applyRBAC()

    alt dbExecutor == nil (dev mode)
        site->>site: skip (pendingHandlers unused)
    else dbExecutor set
        site->>rbac: Init(dbExecutor)
        loop pendingRoles
            site->>rbac: CreateRole(id, code, name, desc)
        end
        site->>rbac: Register(pendingHandlers...)
        Note over rbac: AllowedRoles() read per handler\nPermissions inserted into DB\nAssigned to matching roles\nCache refreshed
    end

    site->>assetmin: ssrBuild(am)
    Note over assetmin: CSS+JS+SVG bundled\nSSR HTML injected
    site->>assetmin: am.RegisterRoutes(mux)
    site->>crudp: cp.RegisterRoutes(mux)
```

## Key invariants

- `RegisterHandlers` is idempotent per handler (crudp deduplicates routes).
- `registerRBAC` only appends — never calls rbac directly before `applyRBAC`.
- `applyRBAC` is called exactly once per process lifetime.
- Order of `SetDB`, `CreateRole`, `RegisterHandlers` is irrelevant — all must precede `Serve`.

## Tests

| Test | Branch covered |
|------|---------------|
| `TestRegisterFlow_DevMode` | pendingHandlers queued, applyRBAC skips, no panic |
| `TestRegisterFlow_OrderIndependence` | RegisterHandlers before CreateRole works |
| `TestRegisterFlow_EmptyHandlers` | RegisterHandlers() → error "no handlers provided" |
| `TestRegisterFlow_PermissionsSeeded` | after Serve with SetDB, rbac has permissions |
