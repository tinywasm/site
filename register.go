package site

import (
	"github.com/tinywasm/crudp"
	"github.com/tinywasm/fmt"
)

// module wraps a handler for site registration
type module struct {
	handler any
	name    string
}

var (
	registeredModules []*module
	cp                *crudp.CrudP
)

// RegisterHandlers registers all handlers with site and crudp
func RegisterHandlers(handlers ...any) error {
	if cp == nil {
		cp = crudp.New()
		cp.SetDevMode(true) // Default dev mode
	}

	if len(handlers) == 0 {
		return nil
	}

	for _, h := range handlers {
		m := &module{handler: h}
		if named, ok := h.(interface{ HandlerName() string }); ok {
			m.name = named.HandlerName()
		}
		registeredModules = append(registeredModules, m)
	}

	if err := cp.RegisterHandlers(handlers...); err != nil {
		fmt.Println("site: crudp registration error:", err)
		return err
	}
	return nil
}

// GetCrudP returns the internal crudp instance
func GetCrudP() *crudp.CrudP {
	return cp
}

// GetModules returns all registered modules
func GetModules() []*module {
	return registeredModules
}
