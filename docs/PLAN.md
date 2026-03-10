# tinywasm/site — Breaking Changes Plan: Explicit SSR Interface

> **Context:** [Architecture Guide](ARCHITECTURE.md)
> **Scope:** Replace the implicit `dom.Component + AllowedRoles('r')==['*']` SSR trigger with an explicit `SSRView` interface. Export the internal `trackedComponentsProvider` as `ComponentTracker`. These are deliberate breaking changes — the ecosystem is in active development.

---

## Development Rules

- **Standard Library Only:** No external assertion libraries in tests. Use `testing` + `reflect`.
- **No Global State:** All changes must be additive to `module.go` (exported interfaces) and surgical in `ssr.build.go`.
- **Max 500 lines per file.** Do not grow `module.go` beyond this.
- **Testing:** Use `gotest` (not `go test`).
- **Documentation First:** Update `docs/ARCHITECTURE.md` before running `gopush`.
- **Breaking changes are intentional.** Do not add backward-compatibility shims.

---

## Context: What Changes and Why

### Removed (breaking)
| Symbol | File | Reason |
|--------|------|--------|
| `trackedComponentsProvider` (unexported) | `ssr.build.go` | Replaced by exported `ComponentTracker` |
| `titleProvider` (unexported) | `ssr.build.go` | `SSRView.Title()` covers this |
| `isPublicReadable()` | `ssr.build.go` | SSR is now opt-in via `SSRView`, not RBAC-derived |
| Implicit `dom.Component` SSR trigger | `ssr.build.go` | Confusing coupling of RBAC + rendering |

### Added (breaking)
| Symbol | File | Purpose |
|--------|------|---------|
| `SSRView` interface | `module.go` | Explicit opt-in for SSR rendering |
| `ComponentTracker` interface | `module.go` | Explicit sub-component CSS/JS/SVG collection |

---

## Step 1 — Add Interfaces to `module.go`

Append to the existing `module.go` file. Do not remove the existing `Module`, `Parameterized`, or `ModuleLifecycle` interfaces.

```go
// SSRView is the explicit contract for server-side rendered modules.
//
// BREAKING CHANGE: Replaces the implicit dom.Component + AllowedRoles('r')==['*'] pattern.
// A handler that previously relied on being a dom.Component with public read access
// must now implement SSRView explicitly.
//
// Register alongside the data handler via site.RegisterHandlers().
// The site engine calls RenderSSR() once during ssrBuild to collect and inject HTML.
type SSRView interface {
	// HandlerName matches the data module this view belongs to.
	// Used for linking and deduplication.
	HandlerName() string
	// Title is the human-readable module name used in navigation and headings.
	Title() string
	// RenderSSR builds the component tree for server-side rendering.
	// Called once at startup during asset bundling. Must be pure: no I/O, no side effects.
	RenderSSR() dom.Component
}

// ComponentTracker is an optional extension for SSRView handlers that compose
// sub-components requiring CSS/JS/SVG collection.
//
// BREAKING CHANGE: Exported replacement for the internal trackedComponentsProvider.
// Implement this when your SSRView builds a tree of components with their own
// CSSProvider, JSProvider, or IconSvgProvider implementations.
type ComponentTracker interface {
	TrackedComponents() []dom.Component
}
```

---

## Step 2 — Refactor `ssr.build.go`

Replace the entire `ssrBuild` function and remove the three unexported interface types.
The rest of the file (imports, `componentRegistry`, etc.) remains unchanged.

**Remove these unexported types** (they are no longer needed):
```go
// DELETE:
type trackedComponentsProvider interface { ... }
type titleProvider interface { ... }
type accessLevel interface { ... }
```

**Replace `ssrBuild` with:**

```go
// ssrBuild collects assets and HTML from all registered SSRView handlers.
func ssrBuild(am *assetmin.AssetMin) error {
	type entry struct {
		view SSRView
		comp dom.Component
	}

	var entries []entry

	// 1. Collect component trees from all SSRView handlers.
	for _, m := range handler.registeredModules {
		view, ok := m.handler.(SSRView)
		if !ok {
			continue
		}
		comp := view.RenderSSR()
		ssr.componentRegistry.register(comp)

		// Collect explicitly declared sub-components (CSS/JS/SVG).
		if tracker, ok := m.handler.(ComponentTracker); ok {
			for _, c := range tracker.TrackedComponents() {
				ssr.componentRegistry.register(c)
			}
		}

		entries = append(entries, entry{view, comp})
	}

	// 2. Inject collected CSS.
	if css := ssr.componentRegistry.collectCSS(); css != "" {
		am.InjectHTML("<style>\n" + css + "</style>\n")
	}

	// 3. Inject collected JS.
	if js := ssr.componentRegistry.collectJS(); js != "" {
		am.InjectHTML("<script>\n" + js + "</script>\n")
	}

	// 4. Inject collected SVG icons into the global sprite.
	for id, svg := range ssr.componentRegistry.collectIcons() {
		am.InjectSpriteIcon(id, svg)
	}

	// 5. Inject module HTML in registration order.
	for _, e := range entries {
		if html := e.comp.RenderHTML(); html != "" {
			am.InjectHTML(html)
		}
	}

	return nil
}

// isPublicReadable is removed. SSR participation is now declared
// explicitly by implementing SSRView, not derived from RBAC configuration.
```

---

## Step 3 — Update `docs/ARCHITECTURE.md`

In the **Interfaces & Components** section, replace the SSR trigger description:

```markdown
### SSR Rendering (`site.SSRView`) — NEW
Register a view handler alongside the data handler. The site engine collects HTML during `ssrBuild`.
- `site.SSRView`: `HandlerName() string`, `Title() string`, `RenderSSR() dom.Component`
- `site.ComponentTracker` (optional): `TrackedComponents() []dom.Component` — declares sub-components for CSS/JS/SVG collection.

**BREAKING:** The previous implicit pattern (`dom.Component` + `AllowedRoles('r')==['*']`) is removed.
```

Remove the entry for `trackedComponentsProvider` (it was internal and is now `ComponentTracker`).
Remove the mention of `isPublicReadable` / `AllowedRoles` as an SSR trigger.

---

## Step 4 — Tests

Add to the existing test suite (or create `ssr_test.go` if none exists):

1. **`TestSSRView_Registered`**: Register a handler implementing `SSRView`. Call `ssrBuild`. Assert `RenderSSR()` was called and its `RenderHTML()` output appears in injected HTML.
2. **`TestSSRView_ComponentTracker`**: Register an `SSRView` that also implements `ComponentTracker`. Assert sub-components are registered (CSS collected).
3. **`TestSSRView_NotTriggered_Without_Interface`**: Register a plain `dom.Component` handler (old pattern). Assert it does NOT appear in injected HTML (breaking change confirmed).

Run with: `gotest`
