package site

import (
	"github.com/tinywasm/fmt"
)

// module definition moved to site.go

// RegisterHandlers registers all handlers with site and crudp
func RegisterHandlers(handlers ...any) error {
	cp.SetDevMode(true) // Default to dev mode, can be changed later

	if len(handlers) == 0 {
		return nil
	}

	for _, h := range handlers {
		m := &module{handler: h}
		if named, ok := h.(interface{ HandlerName() string }); ok {
			m.name = named.HandlerName()
		}
		handler.registeredModules = append(handler.registeredModules, m)
	}

	if err := cp.RegisterHandlers(handlers...); err != nil {
		fmt.Println("site: crudp registration error:", err)
		return err
	}
	return nil
}

// getModules returns all registered modules
func getModules() []*module {
	return handler.registeredModules
}
