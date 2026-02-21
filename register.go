package site

import (
	"github.com/tinywasm/fmt"
)

// RegisterHandlers registers all handlers with site and crudp
func RegisterHandlers(handlers ...any) error {

	if len(handlers) == 0 {
		return fmt.Err("site: no handlers provided")
	}

	for _, h := range handlers {
		name := ""
		if named, ok := h.(interface{ HandlerName() string }); ok {
			name = named.HandlerName()
		}

		if name == "" {
			continue
		}

		// Register as module if it implements Module interface
		if m, ok := h.(Module); ok {
			registerModule(m)
		}

	}

	if err := handler.cp.RegisterHandlers(handlers...); err != nil {
		fmt.Println("site: crudp registration error:", err)
		return err
	}
	// Seed rbac permissions (backend only, no-op on wasm)
	if err := registerRBAC(handlers...); err != nil {
		fmt.Println("site: rbac registration error:", err)
		return err
	}
	// Register assets (SSR only)
	if err := registerAssets(handlers...); err != nil {
		fmt.Println("site: asset registration error:", err)
		return err
	}
	return nil
}

// getModules returns all registered modules
func getModules() []*registeredModule {
	return handler.registeredModules
}
