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
