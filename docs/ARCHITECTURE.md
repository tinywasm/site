# `tinywasm/site` Architecture & Guide (LLM Context)

Isomorphic Go rendering engine orchestrating routing, asset bundling (CSS/JS/SVG), DOM rendering, and RBAC access control.

## 1. Core Principles & Build Split
- **Isomorphic**: Same structs for backend (`!wasm`) and frontend (`wasm`).
- **SSR-First, SPA-Enabled**: Server renders initial HTML; WASM client hydrates.
- **Security Default**: `SetDB` + `SetUserID` mandatory before `Serve`. Bypass with `APP_ENV=development`.
- **Zero-Config Assets**: Auto-bundles CSS/JS/SVG via `tinywasm/assetmin`.

## 2. Server Setup (Backend `!wasm`)
All configuration must happen before `site.Serve(":8080")`.

```go
// 1. Mandatory Security Setup
site.SetDB(&Adp{DB: db}) 
site.SetUserID(func(data ...any) string { 
    if req, ok := data[0].(*http.Request); ok { return req.Header.Get("X-User-ID") }
    return "" 
})

// 2. Queue Roles & Register Handlers (auto-seeds RBAC permissions)
site.CreateRole('a', "Admin", "Desc") 
site.RegisterHandlers(modules.Init()...) 

// 3. Serve (applies RBAC, runs asset bundling, starts server)
site.Serve(":8080") 
```
* **Login Flow**: `site.AssignRole(userID, 'a')` / `site.RevokeRole(userID, 'a')`. Read roles: `site.GetUserRoleCodes(userID)`.
* **Config**: `site.SetCacheSize(3)` (module LRU cache), `site.SetDefaultRoute("home")`.

## 3. Interfaces & Components
A component's capabilities are determined by implementing interfaces (type assertions at registration):

### Module & Navigation Interfaces
Required for a navigable route:
- `site.Module`: `HandlerName() string`, `ModuleTitle() string`, + `dom.Component` (`RenderHTML`, `OnMount`).
Optional:
- `site.Parameterized`: `SetParams(params []string)` (ex: url `#users/123` -> params=`["123"]`).
- `site.ModuleLifecycle`: `BeforeNavigateAway() bool` (block nav if false), `AfterNavigateTo()`.

### Routing & Data (`tinywasm/crudp`)
- `crudp.NamedHandler`: `HandlerName() string`
- `crudp.Creator`, `Reader`, `Updater`, `Deleter`: Maps to POST, GET, PUT, DELETE respectively.
- `crudp.DataValidator`: `ValidateData(action byte, data ...any) error`.

### Access Control & SSR Decision (`crudp.AccessLevel`)
- `AllowedRoles(action byte) []byte` (e.g. action `'r'`, `'c'`, `'u'`, `'d'`).
- **SSR Trigger**: Returning `[]byte{'*'}` for action `'r'` triggers **SSR rendering** (fully indexed HTML).
- **SPA Trigger**: Returning specific roles (e.g., `[]byte{'a'}`) triggers **SPA rendering** (WASM authenticates & renders).

### UI Assets (Backend extraction)
- `site.CSSProvider`: `RenderCSS() string`
- `site.JSProvider`: `RenderJS() string`
- `site.IconSvgProvider`: `IconSvg() map[string]string` (returns map with 1 `"id"` and 1 `"svg"` source).

## 4. Routing & Navigation
- **Data (HTTP)**: Handled by CRUD interfaces mapped to `/{handlerName}/{path...}`.
- **WASM Application Mount**: `site.Mount(parentID string)` initializes the WASM client, mounts the initial module, and blocks forever.
- **WASM SPA Navigation**: `site.Navigate(parentID, "users/123")`. Updates the hashtag to `#users/123` and hydrates state from the LRU cache.

## 5. File Responsibilities (Internal)
* `site.go`: Singleton API delegation.
* `manager.go` / `manager_wasm.go`: Module LRU cache & navigation/hydration.
* `register_ssr.go` (`!wasm`): Asset extraction (CSS/JS/SVG) during `RegisterHandlers`.
* `rbac.back.go` (`!wasm`): Orchestrates `tinywasm/rbac`.

*(For diagrams, see `docs/diagrams/` folder)*
