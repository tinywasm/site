# Site - Orchestration Plan

**Parent Plan**: [API_STANDARDIZATION.md](./API_STANDARDIZATION.md)
**Library**: `github.com/tinywasm/site`
**Status**: Draft - Depends on DOM and Components
**Priority**: 🟡 High (Blocks Full Stack Integration)

---

## Current State Analysis

### What Works ✅
- Hash-based routing (`#home`, `#users`)
- Module registration via `RegisterHandlers()`
- LRU cache (max 3 modules) for performance
- SSR integration with `assetmin`
- CRUDP integration for API routes
- Automatic navigation with `Navigate()`

### What's Unclear ❌
- **Module vs Component distinction**: Both implement `dom.Component`, how do they differ?
- **Navigation placement**: Should navigation be in `site` or `components/nav`?
- **Lifecycle automation**: Is `Unmount → Mount` on navigation automatic or manual?
- **Initial hydration**: How does SSR HTML become interactive?

### Files Overview
```
tinywasm/site/
├── go.mod
├── manager.go         # Navigate(), Start(), module cache
├── module.go          # Module interface definition
├── register.go        # RegisterHandlers()
├── mount.back.go      # SSR: Render(mux), ssrBuild()
├── mount.front.go     # CSR: Mount(), Render()
├── navigation.go      # Hash routing helpers
├── palette.go         # CSS variables (move to components?)
└── example/
    ├── modules/       # Example modules (user, contact)
    └── web/
        ├── server.go  # Backend entry point
        └── client.go  # Frontend entry point
```

---

## Single Responsibility

**Site's ONLY job**: Orchestrate the application - routing, module lifecycle, asset management, and glue between components/modules/backend.

**NOT Site's job**:
- UI primitives (that's `components`)
- DOM manipulation (that's `dom`)
- Business logic (that's `modules`)
- HTTP protocol (that's `crudp`)

**Boundary**: Site is the "conductor" of the orchestra. It tells everyone when to play, but doesn't play the instruments itself.

---

## Module vs Component (Critical Clarification)

### Problem Statement
Both modules and components implement `dom.Component`. What's the difference?

### Definition (Based on Your Feedback)

| Aspect | Component | Module |
|--------|-----------|--------|
| **Purpose** | Reusable UI primitive | Page/feature with business logic |
| **State** | Local UI state only | Domain/business state |
| **Backend** | No backend logic | Has backend logic (CRUD operations) |
| **Location** | `components/` library | `modules/` in user project |
| **Routing** | Not routable | Routable (has a URL route) |
| **Lifecycle** | Managed by parent | Managed by `site` |
| **Example** | Button, Card, Input | UserManagement, ProductCatalog |

### Concrete Distinction

**Component** = Building block, reusable, no backend
```go
// components/button/button.go
type Button struct {
    dom.Component
    Text string
}
// No Read/Create/Update/Delete methods
// No HandlerName() for routing
// Just UI rendering
```

**Module** = Feature/page, routable, has backend
```go
// project/modules/user/user.go
type User struct {
    dom.Component
    ID int
    Name string
}

func (u *User) HandlerName() string { return "users" }
func (u *User) Read(...) any { /* backend logic */ }
// Has CRUD operations
// Registered for routing (#users)
```

### Rule of Thumb

**"Does it need a URL?"**
- ✅ Yes → Module (routable)
- ❌ No → Component (embedded in modules)

**"Does it talk to the database/backend?"**
- ✅ Yes → Module (has business logic)
- ❌ No → Component (pure UI)

---

## Routing Architecture

### Current State (Hash-Based)

```
URL: http://localhost:8080/#users
     ↓
     site.Navigate("users")
     ↓
     Find module with HandlerName() == "users"
     ↓
     Unmount current module
     ↓
     Mount new module
```

**Pros**:
- ✅ Simple (no server-side routing needed)
- ✅ Works with static file serving
- ✅ SPAs are hash-based by default

**Cons**:
- ❌ Ugly URLs (`#` in path)
- ❌ SEO-unfriendly (though less relevant for apps)
- ❌ No nested routes (e.g., `#users/123/edit`)

### Proposed Enhancement (Keep Hash, Add Nesting)

**Support nested routes**:
```
#users          → User list module
#users/123      → User detail module
#users/123/edit → User edit module
```

**Implementation**:
```go
// Parse hash into segments
func parseRoute(hash string) (module string, params []string) {
    parts := strings.Split(hash[1:], "/") // Remove #
    return parts[0], parts[1:]
}

// In Navigate()
module, params := parseRoute(dom.GetHash())
m := findModule(module)
if m != nil {
    m.SetParams(params) // Pass route params to module
    dom.Render("body", m)
}
```

**Decision**: Keep hash-based, add nesting support. History API (pushState) is future enhancement.

---

## Navigation Component

### Question: Where Should Navigation Live?

**Context**: You said navigation should be a component (not in `site`).

**Proposed**: `components/nav/`

```go
// components/nav/nav.go
package nav

type Nav struct {
    components.BaseComponent
    Items []NavItem
}

type NavItem struct {
    Label string
    Route string // e.g., "users", "products"
}

func (n *Nav) Render() dom.Node {
    var items []any
    for _, item := range n.Items {
        items = append(items, dom.Li(
            dom.A(
                dom.Attr("href", "#"+item.Route),
                dom.Text(item.Label),
            ),
        ))
    }
    return dom.Nav(
        dom.Ul(items...),
    )
}
```

**Usage** (in user's main layout):
```go
nav := &nav.Nav{
    Items: []nav.NavItem{
        {Label: "Users", Route: "users"},
        {Label: "Products", Route: "products"},
    },
}
nav.Mount("header")
```

**Why in components?**
- ✅ Navigation IS a UI component (reusable pattern)
- ✅ Keeps `site` focused on orchestration
- ✅ Users can customize/replace navigation easily

**Decision**: Move navigation to `components/nav/`. Remove `navigation.go` from `site`.

---

## Automatic Lifecycle Management

### Current Behavior (manager.go)

```go
func Navigate(parentID string, name string) error {
    // 1. Unmount current module
    if activeModule != nil {
        dom.Unmount(activeModule)
        addToCache(activeModule)
    }

    // 2. Mount new module
    activeModule = target
    dom.SetHash("#" + name)
    return dom.Mount(parentID, target)
}
```

**This is already automatic** ✅ - Good!

**Enhancement**: Add navigation lifecycle hooks for user code:

```go
// In user's module
func (u *User) BeforeNavigateAway() bool {
    // Return false to cancel navigation (e.g., unsaved changes)
    return true
}

func (u *User) AfterNavigateTo() {
    // Called after mount (e.g., fetch data)
}
```

**Implementation**:
```go
// In Navigate()
if bn, ok := activeModule.(BeforeNavigator); ok {
    if !bn.BeforeNavigateAway() {
        return nil // Cancel navigation
    }
}

// ... mount new module ...

if an, ok := target.(AfterNavigator); ok {
    an.AfterNavigateTo()
}
```

**Decision**: Add optional hooks. Automatic lifecycle stays automatic (good default).

---

## SSR to CSR Hydration

### Current Flow

**Server** (mount.back.go):
1. `site.RegisterHandlers(modules.Init()...)` - Register all modules
2. `site.Mount(mux)` - Calls `Render(mux)`
3. `Render()` collects CSS/JS/Icons via `ssrBuild()`
4. `assetmin` serves bundled assets
5. Server renders initial HTML with module content

**Client** (mount.front.go):
1. `site.RegisterHandlers(modules.Init()...)` - Register same modules
2. `site.Mount()` - Calls `Render()` (different implementation)
3. `Render()` calls `Start()` which calls `dom.Hydrate()`
4. `dom.Hydrate()` calls `OnMount()` without re-rendering HTML

**Problem**: How does the client know which module is active on page load?

**Solution**: Server embeds initial route in HTML:
```html
<script>
    window.__INITIAL_ROUTE__ = "users"; // or read from hash
</script>
```

**Client reads**:
```go
func Start() error {
    hash := dom.GetHash()
    if hash == "" {
        hash = "#home" // Default
    }

    name := hash[1:] // Remove #
    m := findModule(name)
    return dom.Hydrate("app", m) // Assumes <div id="app"> in HTML
}
```

**Decision**: Use hash from URL (no need for script tag). Client reads `window.location.hash`.

---

## LRU Cache (Performance Optimization)

### Current Implementation (manager.go)

```go
var cache []Module
var maxCache = 3

func addToCache(m Module) {
    // Simple LRU: remove oldest if full
    if len(cache) >= maxCache {
        cache = cache[1:]
    }
    cache = append(cache, m)
}
```

**This is good** ✅ - User confirmed.

**Reasoning**: Users typically navigate between 2-3 modules (e.g., list → detail → list). Caching 3 avoids re-creating components.

**Enhancement**: Make cache size configurable (optional):
```go
func SetCacheSize(size int) {
    maxCache = size
}
```

**Decision**: Keep current implementation, add `SetCacheSize()` for power users.

---

## API Simplification

### Current API (server.go example)

```go
mux := http.NewServeMux()
site.RegisterHandlers(modules.Init()...)
site.Mount(mux)
http.ListenAndServe(":8080", mux)
```

**Problem**: User must create `http.ServeMux` and `ListenAndServe` manually.

### Proposed API (Simplified)

```go
site.RegisterHandlers(modules.Init()...)
site.Serve(":8080") // One-liner
```

**Implementation**:
```go
//go:build !wasm

func Serve(addr string) error {
    mux := http.NewServeMux()
    if err := Mount(mux); err != nil {
        return err
    }
    return http.ListenAndServe(addr, mux)
}
```

**Pros**:
- ✅ 1 line instead of 3
- ✅ Still allows custom mux if needed (via `Mount(customMux)`)

**Decision**: Add `Serve()` helper, keep `Mount()` for advanced use.

---

## Client API (Matching Simplicity)

### Current API (client.go example)

```go
site.RegisterHandlers(modules.Init()...)
site.Mount()
select {}
```

**This is already simple** ✅

**Enhancement**: Remove `select {}` requirement by handling it internally:

```go
func Mount() error {
    if err := Render(); err != nil {
        return err
    }
    select {} // Block forever (WASM apps don't exit)
}
```

**Usage**:
```go
func main() {
    site.RegisterHandlers(modules.Init()...)
    site.Mount() // Blocks automatically
}
```

**Decision**: Make `Mount()` block automatically. Developers shouldn't need `select {}`.

---

## Module Interface (Final Definition)

### Current Interface (module.go)

```go
type Module interface {
    dom.Component
    HandlerName() string
    ModuleTitle() string
}
```

**Enhancement**: Add optional lifecycle hooks:

```go
type Module interface {
    dom.Component
    HandlerName() string  // Required: Route name
    ModuleTitle() string  // Required: Display name
}

// Optional lifecycle hooks
type ModuleLifecycle interface {
    BeforeNavigateAway() bool // Return false to cancel
    AfterNavigateTo()          // Called after mount
}

// Optional route params
type Parameterized interface {
    SetParams(params []string)
}
```

**Example**:
```go
type User struct {
    dom.BaseComponent
    ID int
}

func (u *User) HandlerName() string { return "users" }
func (u *User) ModuleTitle() string { return "User Management" }

func (u *User) SetParams(params []string) {
    if len(params) > 0 {
        u.ID = parseInt(params[0]) // Route: #users/123
    }
}

func (u *User) AfterNavigateTo() {
    // Fetch user data based on u.ID
}
```

---

## Integration with Components

### How Modules Use Components

**Pattern**: Modules compose components in their `Render()` method:

```go
//go:build wasm

func (u *User) Render() dom.Node {
    return dom.Div(
        &nav.Nav{Items: u.navItems}, // Component
        dom.H1(dom.Text("Users")),
        &button.Button{Text: "Add", Variant: "primary"}, // Component
        &table.Table{Data: u.users}, // Component
    )
}
```

**Key Point**: Modules are responsible for layout and composition. Components provide the building blocks.

**Decision**: No changes needed. DOM's DSL already supports nesting components (from uncommitted changes).

---

## Asset Management (assetmin Integration)

### Current Flow (mount.back.go)

```go
func Render(mux *http.ServeMux) error {
    am := assetmin.NewAssetMin(&assetmin.Config{
        OutputDir: "./public",
        DevMode: handler.DevMode,
    })

    ssrBuild(am) // Collect CSS/Icons from modules
    am.RegisterRoutes(mux) // Serve /public/style.css, /public/sprite.svg
    handler.cp.RegisterRoutes(mux) // CRUDP API routes
}
```

**This is good** ✅

**Enhancement**: Allow custom output directory:
```go
func SetOutputDir(dir string) {
    // Default: "./public"
}
```

**Decision**: Add configuration helpers, keep current architecture.

---

## API Summary (Final)

### Server API

```go
package site

// Registration
func RegisterHandlers(handlers ...any) error

// Serving (simple)
func Serve(addr string) error

// Serving (advanced, custom mux)
func Mount(mux *http.ServeMux) error

// Configuration
func SetDevMode(enabled bool)
func SetOutputDir(dir string)
func SetCacheSize(size int)
func SetUserRoles(fn func(...any) []byte)
```

### Client API

```go
package site

// Registration (same as server)
func RegisterHandlers(handlers ...any) error

// Mounting (blocks automatically)
func Mount() error

// Navigation (usually automatic via hash changes)
func Navigate(moduleName string) error

// Rendering (called by Mount, rarely used directly)
func Render() error
```

### Module Interface

```go
package site

type Module interface {
    dom.Component        // ID(), SetID(), RenderHTML() or Render()
    HandlerName() string // Route name (e.g., "users")
    ModuleTitle() string // Display name (e.g., "User Management")
}

// Optional interfaces
type ModuleLifecycle interface {
    BeforeNavigateAway() bool
    AfterNavigateTo()
}

type Parameterized interface {
    SetParams(params []string)
}
```

---

## Example Usage (Full Stack)

### Server (web/server.go)

```go
//go:build !wasm

package main

import (
    "myapp/modules"
    "github.com/tinywasm/site"
)

func main() {
    site.RegisterHandlers(modules.Init()...)
    site.Serve(":8080") // One line!
}
```

### Client (web/client.go)

```go
//go:build wasm

package main

import (
    "myapp/modules"
    "github.com/tinywasm/site"
)

func main() {
    site.RegisterHandlers(modules.Init()...)
    site.Mount() // Blocks automatically
}
```

### Module (modules/user/user.go)

```go
package user

import (
    "github.com/tinywasm/dom"
    "github.com/tinywasm/components/button"
    "github.com/tinywasm/components/table"
)

type User struct {
    dom.BaseComponent
    users []UserData
}

func (u *User) HandlerName() string { return "users" }
func (u *User) ModuleTitle() string { return "User Management" }

func (u *User) Render() dom.Node {
    return dom.Div(
        dom.H1(dom.Text("Users")),
        &button.Button{Text: "Add User", Variant: "primary"},
        &table.Table{Data: u.users},
    )
}

func (u *User) AfterNavigateTo() {
    // Fetch users via CRUDP
    u.users = fetchUsers()
    u.Render() // Re-render with data
}
```

**Lines of code**: ~15 (vs ~40 in current API)

---

## Open Questions for Approval

### ❓ Q1: Navigation Component
**Move navigation to `components/nav/` (remove from `site`)?**
- [ ] Yes, navigation is a component
- [ ] No, keep it in site
- [ ] Other suggestion: _______________

### ❓ Q2: Nested Routes
**Support nested routes (`#users/123/edit`)?**
- [ ] Yes, add nesting support
- [ ] No, keep flat routes only
- [ ] Defer to later

### ❓ Q3: Lifecycle Hooks
**Add optional `BeforeNavigateAway()` / `AfterNavigateTo()` hooks?**
- [ ] Yes, add them (optional interfaces)
- [ ] No, keep current lifecycle only

### ❓ Q4: API Simplification
**Add `site.Serve()` one-liner for server?**
- [ ] Yes, add it
- [ ] No, keep current `Mount(mux)` only

**Make client `Mount()` block automatically (remove `select {}`)?**
- [ ] Yes, block automatically
- [ ] No, keep `select {}` in user code

---

## Next Steps After Approval

1. Implement nested route parsing
2. Add lifecycle hooks (optional interfaces)
3. Move navigation to `components/nav/`
4. Add `Serve()` helper for server
5. Make client `Mount()` block automatically
6. Update examples and documentation

**Estimated Effort**: 2-3 days

---

## Relationship to Other Plans

**Depends On**:
- ✅ DOM API (for `dom.Render`, `dom.Hydrate`, lifecycle hooks)
- ✅ Components (for `components/nav/`)

**Enables**:
- ✅ Full-stack application development
- ✅ Module authors can use standardized APIs

**Blocks**:
- ❌ Nothing (site is final orchestration layer)

---

**Ready for approval?** This completes the core API standardization. The last remaining plan is a brief spec for `modules/` (README only, no implementation).
