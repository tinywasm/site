# Access Control Flow

> **Status:** Current — February 2026

Request-time authorization check. Wired at `Mount` time via `crudp.SetAccessCheck`.
`getUserID` is either provided by `site.SetUserID` or auto-wired from session cookie
when `tinywasm/user` is initialized.

```mermaid
sequenceDiagram
    participant Client
    participant crudp
    participant site
    participant rbac

    Client->>crudp: HTTP request (resource, action)
    crudp->>site: SetAccessCheck closure(resource, action, *http.Request)

    site->>site: getUserID(*http.Request)
    alt getUserID == nil (SetUserID not configured)
        site-->>crudp: false (deny)
        crudp-->>Client: 403 Forbidden
    else getUserID configured
        site->>site: extract userID from request
        alt userID == "" (anonymous)
            site-->>crudp: false (deny)
            crudp-->>Client: 401 Unauthorized
        else userID present
            site->>rbac: HasPermission(userID, resource, action)
            rbac->>rbac: cache lookup (userID → roles → permissions)
            alt cache hit
                rbac-->>site: (allowed bool, nil)
            else cache miss
                rbac->>rbac: DB query → warm cache entry
                rbac-->>site: (allowed bool, nil)
            end
            alt allowed == true
                site-->>crudp: true
                crudp-->>Client: 200 OK + payload
            else allowed == false
                site-->>crudp: false
                crudp-->>Client: 403 Forbidden
            end
        end
    end
```

## Notes

- `AllowedRoles('r') == []byte{'*'}` (public read) → `HasPermission` still called but
  the wildcard permission is seeded with `*` role, matching all users including anonymous.
- Cache is in-memory (`sync.RWMutex`), populated at `rbac.Init` and updated on `AssignRole`/`RevokeRole`.
- Zero DB I/O on the hot path — all reads served from cache.

## Tests

| Test | Branch covered |
|------|---------------|
| `TestAccessCheck_Anonymous` | empty userID → deny |
| `TestAccessCheck_NoUserIDFunc` | getUserID == nil → deny |
| `TestAccessCheck_Allowed` | valid userID + matching role → allow |
| `TestAccessCheck_Denied` | valid userID + no matching role → deny |
| `TestAccessCheck_PublicRead` | AllowedRoles('r')=='*' → allow anonymous |
