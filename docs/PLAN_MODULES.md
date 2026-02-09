# Plan: Site Module System & Router

This plan details the implementation of the Module System, Router, and Hydration integration within `tinywasm/site`.

## Goal
Implement a centralized module manager that handles routing, module switching (mount/unmount), and hydration of the initial state, using a linear slice approach for performance in TinyGo.

## Proposed Changes

### [tinywasm/site]

#### [NEW] [module.go](file:///home/cesar/Dev/Pkg/tinywasm/site/module.go)
- Define `Module` interface:
    ```go
    type Module interface {
        dom.Component
        HandlerName() string
        ModuleTitle() string
    }
    ```

#### [NEW] [manager.go](file:///home/cesar/Dev/Pkg/tinywasm/site/manager.go)
- Implement `site.Register(m Module)`.
- Implement `site.Start()` (WASM Entrypoint):
    - Parse URL hash to find initial module.
    - Call `dom.Hydrate("root", modules[i])`.
- Implement `site.Navigate(name string)`:
    - Look up module by `HandlerName()`.
    - `dom.Unmount(activeModule)`.
    - Update `activeModule`.
    - `dom.Render("root", newModule)`.
    - Update URL hash.

#### [NEW] [manager_test.go](file:///home/cesar/Dev/Pkg/tinywasm/site/manager_test.go)
- Unit tests for logic (independent of DOM/WASM where possible, or using mocks).

## Verification Plan

### Automated Tests
- Create `manager_test.go`:
    - **Test 1:** `TestRegisterAndFind`: Register modules and verify linear search finds them.
    - **Test 2:** `TestNavigateLogic`: Verify `Navigate` updates the `activeModule` pointer (mocking DOM calls if necessary or separating logic from DOM side effects).

### Manual Verification
- **Browser Test:**
    1.  Compile and run the `site` example.
    2.  Land on Homepage -> Verify it's hydrated (click works immediately).
    3.  Click "User" link -> Verify "Home" unmounts, "User" mounts.
    4.  Verify URL hash changes.
    5.  Refresh -> Verify "User" hydrates correctly.
