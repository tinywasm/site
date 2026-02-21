# Setup Flow

> **Status:** Current — February 2026
> `applyUser()` marked as **[planned]** — implemented in a subsequent plan.

Two initialization paths depending on security configuration.

```mermaid
sequenceDiagram
    participant App as Application
    participant site
    participant rbac
    participant user as user [planned]

    App->>site: SetDB(adapter) [optional]
    App->>site: SetUserID(fn) [optional, manual override]
    App->>site: CreateRole(code, name, desc) [queued]
    App->>site: RegisterHandlers(hs...)
    Note over site: pendingHandlers queued
    App->>site: Serve(":8080")
    site->>site: Mount(mux)

    alt APP_ENV=development (or -dev flag)
        site->>site: applyRBAC() → dbExecutor==nil → skip
        site->>site: applyUser() → dbExecutor==nil → skip [planned]
        site->>site: DevMode=true → skip security validation
        site-->>App: ready (no DB required)
    else Production (SetDB called)
        site->>rbac: Init(dbExecutor)
        rbac->>rbac: runMigrations()
        rbac->>rbac: warmCache()
        rbac-->>site: ok
        site->>site: rbacInitialized = true
        site->>site: wire SetAccessCheck closure
        loop pendingRoles
            site->>rbac: CreateRole(id, code, name, desc)
        end
        site->>rbac: Register(pendingHandlers...)
        site->>user: Init(dbExecutor) [planned]
        user->>user: runMigrations() [planned]
        user-->>site: ok [planned]
        site->>site: userInitialized = true [planned]
        alt getUserID == nil
            site->>site: auto-wire getUserID from session cookie [planned]
        end
        alt getUserID still nil
            site-->>App: error: SetUserID required
        else getUserID set
            site-->>App: ready
        end
    end
```

## Tests

| Test | Branch covered |
|------|---------------|
| `TestDevMode_BypassesSecurityCheck` | `APP_ENV=development` → Mount succeeds without SetDB |
| `TestDevMode_FlagArg` | `-dev` flag → DevMode=true |
| `TestProduction_RequiresSetDB` | no SetDB, no dev mode → Mount returns error |
| `TestProduction_SetDB_Success` | SetDB + SetUserID → Mount succeeds |
| `TestProduction_SetDB_NoUserID` | SetDB without SetUserID (no user module) → Mount returns error |
