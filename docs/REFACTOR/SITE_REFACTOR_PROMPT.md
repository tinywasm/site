# Site Refactor - Execution Prompt

**Library**: `github.com/tinywasm/site`
**Location**: ``
**Status**: Ready to execute (after DOM + Components)
**Estimated Time**: 2-3 days
**Priority**: 🟡 High (depends on DOM and Components completion)

---

## Context

You are refactoring the `tinywasm/site` library to:
1. **Simplify server/client APIs** (one-liners)
2. **Add nested route support** (`#users/123/edit`)
3. **Remove navigation from site** (moved to `components/nav`)
4. **Add optional lifecycle hooks** (BeforeNavigateAway, AfterNavigateTo)
5. **Integrate with refactored DOM and Components**

**Goal**: Simplify application orchestration, reduce boilerplate by 50%.

---

## Prerequisites

### ✅ Verify Dependencies Complete
Before starting, verify:
- [ ] **DOM refactor complete**: Fluent API, Elm pattern working
- [ ] **Components refactor complete**: 7 components exist, nav component ready
- [ ] Both pass `gotest`

**If not complete**: Coordinate with DOM/Components agents or wait.

---

## Current State

### Files to Review
- `manager.go` - Navigate(), LRU cache
- `module.go` - Module interface
- `register.go` - RegisterHandlers()
- `mount.back.go` - SSR: Render(mux)
- `mount.front.go` - CSR: Mount(), Render()
- `navigation.go` - Hash routing (to be removed)
- `example/` - Current examples

### Current Issues
- ❌ Verbose server setup (3 lines instead of 1)
- ❌ Client requires manual `select {}` blocking
- ❌ Navigation code mixed with site orchestration
- ❌ No nested route support
- ❌ No lifecycle hooks for navigation

---

## Approved Decisions

### Decision 1: API Simplification

**Server** (before):
```go
mux := http.NewServeMux()
site.RegisterHandlers(modules.Init()...)
site.Mount(mux)
http.ListenAndServe(":8080", mux)
```

**Server** (after):
```go
site.RegisterHandlers(modules.Init()...)
site.Serve(":8080") // One-liner
```

**Client** (before):
```go
site.RegisterHandlers(modules.Init()...)
site.Mount()
select {} // Manual blocking
```

**Client** (after):
```go
site.RegisterHandlers(modules.Init()...)
site.Mount() // Blocks automatically
```

**Requirements**:
- Add `site.Serve(addr)` helper for server
- Make client `Mount()` block internally
- Keep `Mount(mux)` for advanced use cases

---

### Decision 2: Nested Route Support

**Current**: `#users` (flat)
**New**: `#users/123/edit` (nested)

**Parsing**:
```go
// #users/123/edit → module="users", params=["123", "edit"]
func parseRoute(hash string) (module string, params []string) {
    parts := strings.Split(hash[1:], "/") // Remove #
    if len(parts) == 0 {
        return "", nil
    }
    return parts[0], parts[1:]
}
```

**Module interface** (add optional):
```go
type Parameterized interface {
    SetParams(params []string)
}

// Usage in module
func (u *User) SetParams(params []string) {
    if len(params) > 0 {
        u.ID = parseInt(params[0]) // Route: #users/123
    }
}
```

**Requirements**:
- Parse hash into module name + params
- Pass params to module via `SetParams()` if implemented
- Backward compatible (modules without params still work)

---

### Decision 3: Navigation Moved to Components

**Remove**: `site/navigation.go`
**Reason**: Navigation is a UI component, belongs in `components/nav`

**New flow**:
1. User creates `nav.Nav` component (from `components/nav`)
2. Renders it in their layout
3. Navigation links use `href="#route"`
4. Site's hash change listener triggers module switch

**Site still handles**:
- Hash change detection (`OnHashChange`)
- Module switching (`Navigate`)
- LRU cache management

**Site does NOT handle**:
- Navigation UI rendering (that's `components/nav`)

---

### Decision 4: Lifecycle Hooks

**Add optional interfaces**:
```go
type ModuleLifecycle interface {
    BeforeNavigateAway() bool  // Return false to cancel
    AfterNavigateTo()          // Called after mount
}
```

**Example usage**:
```go
func (u *User) BeforeNavigateAway() bool {
    if u.hasUnsavedChanges {
        return confirm("Unsaved changes. Leave anyway?")
    }
    return true
}

func (u *User) AfterNavigateTo() {
    // Fetch user data
    u.users = u.fetchUsers()
    u.Update()
}
```

**Requirements**:
- Check interface before navigation
- Call hooks at appropriate times
- Optional (modules don't need to implement)

---

### Decision 5: Keep LRU Cache

**Current**: Max 3 modules cached
**Decision**: Keep as-is (good balance)

**Optional**: Add `SetCacheSize(size int)` for power users

---

## Implementation Tasks

### Task 1: Simplify Server API

**File**: `mount.back.go` (or new `serve.go`)

**Add `Serve()` helper**:
```go
//go:build !wasm

package site

import "net/http"

// Serve starts the server on the given address (one-liner helper)
func Serve(addr string) error {
    mux := http.NewServeMux()
    if err := Mount(mux); err != nil {
        return err
    }
    return http.ListenAndServe(addr, mux)
}

// Mount remains for advanced use (custom mux, middleware, etc.)
func Mount(mux *http.ServeMux) error {
    // Existing implementation (rename from Render if needed)
    // ...
}
```

**Keep backward compatibility**: `Mount(mux)` still works for custom setups.

---

### Task 2: Simplify Client API

**File**: `mount.front.go`

**Make `Mount()` block automatically**:
```go
//go:build wasm

package site

// Mount hydrates the initial module and blocks forever
func Mount() error {
    if err := Render(); err != nil {
        return err
    }
    select {} // Block automatically (WASM apps don't exit)
}

// Render remains for advanced use (testing, custom lifecycle)
func Render() error {
    return Start("app") // Or read from config
}
```

**Alternate approach** (if you want `Mount()` to take parent ID):
```go
func Mount(parentID string) error {
    if err := Start(parentID); err != nil {
        return err
    }
    select {}
}
```

---

### Task 3: Implement Nested Routes

**File**: `manager.go`

**Add route parsing**:
```go
import "strings"

// parseRoute extracts module name and params from hash
func parseRoute(hash string) (module string, params []string) {
    if hash == "" || hash == "#" {
        return "home", nil // Default route
    }

    // Remove # and split by /
    parts := strings.Split(hash[1:], "/")
    if len(parts) == 0 {
        return "home", nil
    }

    return parts[0], parts[1:]
}
```

**Update `Navigate()` to use params**:
```go
func Navigate(parentID string, hash string) error {
    moduleName, params := parseRoute(hash)

    if activeModule != nil && activeModule.HandlerName() == moduleName {
        // Same module, just update params
        if p, ok := activeModule.(Parameterized); ok {
            p.SetParams(params)
        }
        return nil
    }

    target := findModule(moduleName)
    if target == nil {
        return nil // Or 404
    }

    // 1. Check if current module allows navigation away
    if activeModule != nil {
        if lc, ok := activeModule.(ModuleLifecycle); ok {
            if !lc.BeforeNavigateAway() {
                return nil // Cancelled
            }
        }
        dom.Unmount(activeModule)
        addToCache(activeModule)
    }

    // 2. Set params on new module
    if p, ok := target.(Parameterized); ok {
        p.SetParams(params)
    }

    // 3. Mount new module
    activeModule = target
    dom.SetHash(hash)
    if err := dom.Mount(parentID, target); err != nil {
        return err
    }

    // 4. Call AfterNavigateTo hook
    if lc, ok := target.(ModuleLifecycle); ok {
        lc.AfterNavigateTo()
    }

    return nil
}
```

**Update `Start()` to use route parsing**:
```go
func Start(parentID string) error {
    hash := dom.GetHash()
    moduleName, params := parseRoute(hash)

    m := findModule(moduleName)
    if m == nil {
        return fmt.Errf("module not found: %s", moduleName)
    }

    // Set params
    if p, ok := m.(Parameterized); ok {
        p.SetParams(params)
    }

    activeModule = m
    return dom.Hydrate(parentID, m)
}
```

---

### Task 4: Add Lifecycle Interfaces

**File**: `module.go`

**Add optional interfaces**:
```go
package site

import "github.com/tinywasm/dom"

// Module represents a top-level site module (routable page)
type Module interface {
    dom.Component
    HandlerName() string  // Route name (e.g., "users")
    ModuleTitle() string  // Display name (e.g., "User Management")
}

// Parameterized modules can receive route parameters
type Parameterized interface {
    SetParams(params []string)
}

// ModuleLifecycle provides hooks for navigation events
type ModuleLifecycle interface {
    BeforeNavigateAway() bool // Return false to cancel navigation
    AfterNavigateTo()         // Called after module is mounted
}
```

---

### Task 5: Remove Navigation Code

**File**: `navigation.go`

**Delete this file** or mark as deprecated:
```go
// DEPRECATED: Navigation UI moved to components/nav
// This file is kept for backward compatibility only.
// Use github.com/tinywasm/components/nav instead.

package site

// OnHashChange, GetHash, SetHash moved to dom package
// Redirect to dom.OnHashChange, etc.
```

**Update other files** that import from `navigation.go`:
- Replace `site.OnHashChange` with `dom.OnHashChange`
- Replace `site.GetHash` with `dom.GetHash`
- Replace `site.SetHash` with `dom.SetHash`

---

### Task 6: Add Configuration Helpers

**File**: `config.go` (new)

```go
package site

var (
    config = &Config{
        CacheSize: 3,
        DefaultRoute: "home",
    }
)

type Config struct {
    CacheSize    int
    DefaultRoute string
}

// SetCacheSize configures module cache size (default: 3)
func SetCacheSize(size int) {
    config.CacheSize = size
    maxCache = size
}

// SetDefaultRoute configures default route (default: "home")
func SetDefaultRoute(route string) {
    config.DefaultRoute = route
}
```

---

### Task 7: Update Examples

**Server** (`example/web/server.go`):
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

**Client** (`example/web/client.go`):
```go
//go:build wasm

package main

import (
    "myapp/modules"
    "github.com/tinywasm/site"
)

func main() {
    site.RegisterHandlers(modules.Init()...)
    site.Mount("app") // Blocks automatically, no select{}
}
```

**Module with params** (`example/modules/user/user.go`):
```go
package user

import "github.com/tinywasm/dom"

type User struct {
    dom.BaseComponent
    ID    int
    users []UserData
}

func (u *User) HandlerName() string { return "users" }
func (u *User) ModuleTitle() string { return "User Management" }

// Handle route params (#users/123)
func (u *User) SetParams(params []string) {
    if len(params) > 0 {
        u.ID = parseInt(params[0])
    }
}

// Lifecycle hooks
func (u *User) AfterNavigateTo() {
    u.users = u.fetchUsers()
    u.Update()
}

func (u *User) BeforeNavigateAway() bool {
    if u.hasUnsavedChanges {
        return confirm("Unsaved changes. Continue?")
    }
    return true
}

func (u *User) Render() dom.Node {
    return dom.Div().
        Class("user-module").
        Append(dom.H1().Text("Users")).
        Append(/* user list */).
        ToNode()
}
```

---

### Task 8: Update Tests

**Add tests for nested routes**:

`route_parsing_test.go`:
```go
func TestParseRoute(t *testing.T) {
    tests := []struct{
        hash   string
        module string
        params []string
    }{
        {"#users", "users", nil},
        {"#users/123", "users", []string{"123"}},
        {"#users/123/edit", "users", []string{"123", "edit"}},
        {"", "home", nil},
    }

    for _, tt := range tests {
        mod, params := parseRoute(tt.hash)
        if mod != tt.module {
            t.Errorf("expected module %s, got %s", tt.module, mod)
        }
        // Compare params...
    }
}
```

`lifecycle_test.go`:
```go
func TestModuleLifecycle(t *testing.T) {
    type TestModule struct {
        dom.BaseComponent
        beforeCalled bool
        afterCalled  bool
    }

    func (m *TestModule) BeforeNavigateAway() bool {
        m.beforeCalled = true
        return true
    }

    func (m *TestModule) AfterNavigateTo() {
        m.afterCalled = true
    }

    // Test hooks are called
}
```

---

## Success Criteria

Before marking as complete, verify:

### ✅ Functionality
- [ ] `site.Serve(":8080")` works (server)
- [ ] `site.Mount("app")` blocks automatically (client)
- [ ] Nested routes work (`#users/123/edit`)
- [ ] `SetParams()` called with route params
- [ ] Lifecycle hooks called at correct times
- [ ] LRU cache still works
- [ ] Navigation moved to `components/nav`

### ✅ Testing
- [ ] Route parsing tests pass
- [ ] Lifecycle tests pass
- [ ] Integration test (server + client)
- [ ] Run `gotest` in site directory

### ✅ Examples
- [ ] Server example uses new API
- [ ] Client example uses new API
- [ ] Module example shows params and hooks
- [ ] Examples compile with TinyGo

### ✅ Documentation
- [ ] README.md updated with new API
- [ ] Migration guide created (old API → new API)
- [ ] Module interface documented
- [ ] Lifecycle hooks documented

---

## Files to Modify/Create

| File | Action | Description |
|------|--------|-------------|
| `serve.go` | Create | Add Serve() helper |
| `mount.front.go` | Modify | Make Mount() block automatically |
| `manager.go` | Modify | Add nested route support, lifecycle hooks |
| `module.go` | Modify | Add Parameterized, ModuleLifecycle interfaces |
| `navigation.go` | Delete/Deprecate | Move to components/nav |
| `config.go` | Create | Configuration helpers |
| `example/web/server.go` | Update | Use new API |
| `example/web/client.go` | Update | Use new API |
| `example/modules/user/` | Update | Show params + hooks |
| `route_parsing_test.go` | Create | Test route parsing |
| `lifecycle_test.go` | Create | Test lifecycle hooks |
| `README.md` | Update | Document new API |
| `MIGRATION.md` | Create | Migration guide |

---

## Implementation Order

1. **Add lifecycle interfaces** (foundation)
2. **Implement route parsing** (parseRoute function)
3. **Update Navigate() with params and hooks** (core logic)
4. **Add Serve() helper** (server API)
5. **Make Mount() block** (client API)
6. **Remove/deprecate navigation.go** (cleanup)
7. **Add config helpers** (utilities)
8. **Update examples** (documentation)
9. **Write tests** (validation)
10. **Update docs** (guides)

---

## Migration Guide (Create as MIGRATION.md)

```markdown
# Site API Migration Guide

## Server Changes

### Before
```go
mux := http.NewServeMux()
site.RegisterHandlers(modules.Init()...)
site.Mount(mux)
http.ListenAndServe(":8080", mux)
```

### After
```go
site.RegisterHandlers(modules.Init()...)
site.Serve(":8080")
```

**Note**: If you need custom middleware, use `site.Mount(mux)` instead.

## Client Changes

### Before
```go
site.RegisterHandlers(modules.Init()...)
site.Mount()
select {}
```

### After
```go
site.RegisterHandlers(modules.Init()...)
site.Mount("app")
```

**Note**: `select {}` is now automatic.

## Navigation Changes

### Before
Navigation UI was in `site/navigation.go`.

### After
Use `components/nav`:

```go
import "github.com/tinywasm/components/nav"

navigation := &nav.Nav{
    Items: []nav.NavItem{
        {Label: "Users", Route: "users"},
        {Label: "Products", Route: "products"},
    },
}
dom.Render("header", navigation)
```

## Nested Routes (New Feature)

### Before
Only flat routes: `#users`

### After
Nested routes: `#users/123/edit`

```go
func (u *User) SetParams(params []string) {
    if len(params) > 0 {
        u.ID = parseInt(params[0])
    }
}
```

## Lifecycle Hooks (New Feature)

```go
func (u *User) BeforeNavigateAway() bool {
    return !u.hasUnsavedChanges // Cancel if unsaved
}

func (u *User) AfterNavigateTo() {
    u.fetchData()
    u.Update()
}
```

---

## Questions/Ambiguities?

If you encounter decisions not covered here:
1. **Read** [SITE_ORCHESTRATION.md](./SITE_ORCHESTRATION.md) for full context
2. **Follow principles**: Simplify API, no magic, explicit > implicit
3. **Coordinate**: If navigation component not ready, wait or ask Components agent

---

## Completion

When done:
1. Commit: `refactor(site): simplify API, add nested routes, add lifecycle hooks`
2. Run `gotest` and paste output
3. Test example: `cd example && gotest`
4. Create MIGRATION.md guide
5. Report completion to coordinator

---

**Status**: Ready to execute after DOM and Components refactors complete.
