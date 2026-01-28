# Prompt: Initial Architecture for Site

## Context
`site` is a new package designed to be the administrator/manager for websites built with `tinydom`. It serves as the bridge between the backend (SSR, asset management) and the frontend (WASM, SPA routing).

## Objectives

### 1. Core Structure
**Requirement**:
*   Implement the `Site` struct in `site.go`.
*   It should act as the central registry for the application's pages.

### 2. Module & Template Architecture
**Requirement**:
*   **Modules as Pages**: Each module (e.g., `modules/users`, `modules/contact`) represents a functional part of the site, typically a Page.
*   **Templates**: Modules should NOT build their UI from scratch in Go code alone. They should use templates (e.g., stored in `web/template`) to define their structure.
*   **Page Interface**: Modules must provide metadata and methods to help build the site structure.
    *   `PageName() string`: Returns the name/title of the page.
    *   `IconSvg() string`: Returns the raw SVG string for the module's icon (used in navigation).
    *   `RenderHTML() string`: (via `tinydom` component) Returns the initial HTML content using the pre-configured templates.
*   **Reusability**: Templates should be generic enough to be reused across modules, changing only specific values (titles, content slots).

### 3. Router Construction
**Requirement**:
*   `site` should aggregate these pages to build the application router.
*   Implement `AddPages(pages ...any)` (or a specific interface type like `...Page`).
    *   *Note*: `AddPages` is preferred over `AddModules` as it clearly indicates we are adding navigable pages to the site.
*   Internally, `Site` iterates over these pages to:
    1.  Register routes (e.g., `/users`, `/contact`).
    2.  build the navigation menu using `PageName()` and `IconSvg()`.
    3.  Map routes to their respective components for rendering (SSR and Client-side).

### 4. Integration with TinyDOM
**Requirement**:
*   `site` orchestrates `tinydom` components.
*   It ensures that when a route is requested, the correct component's `RenderHTML()` is called (for SSR) or mounted (for WASM).

## Usage Example (Conceptual)

```go
// In modules/modules.go
// Init returns all the pages that make up the site
func Init() []any {
    return []any{
        users.Add(),
        contact.Add(),
        // ... other pages
    }
}

// In web/site/site.go

func Init(logger func(msg ...any)) *site.Site {

    config := site.configSite{
        Logger: logger,
        ColorScheme: site.ColorScheme{
            // Brand Colors
            Primary:       "#000000", // Main brand color (buttons, links)
            Secondary:     "#000000", // Accent color
            // Backgrounds for Light Mode or Black Mode
            Background:    "#FFFFFF", // App/Page background
            Surface:       "#F5F5F5", // Card/Modal/Sidebar background (elevation)
            // Content
            TextPrimary:   "#000000", // High emphasis text
            TextSecondary: "#666666", // Medium emphasis text (subtitles)
            TextOnPrimary: "#FFFFFF", // Text color on top of Primary elements
            // UI
            Border:        "#E0E0E0", // Dividers, borders
            // Status
            Error:         "#FF0000",
            Success:       "#00FF00",
        },
    }

    newSite := site.New(config)

    // Register pages using the centralized Init function
    newSite.AddPages(modules.Init()...)

    return newSite
}
// WebSite is now ready to handle routing and rendering

// In web/server.go (go:build !wasm)
    logger := func(msg ...any) {
        log.Println(msg...)
    }
    ssr := site.Init(logger)

    http.HandleFunc("/", ssr.HandleRequest)
---
//In web/client.go (go:build wasm)
    logger := func(msg ...any) {
        js.Global().Get("console").Call("log", msg...)
    }
    spa := site.Init(logger)

    spa.Mount()
 select{} //for WASM to keep running
```

## Deliverables
*   Define the `Site` struct and `AddPages` method.
*   Define the interface that pages must implement (Name, Icon, Component methods).
*   Implement the logic to build the router and navigation from these pages.
