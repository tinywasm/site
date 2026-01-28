package site

import (
	"github.com/tinywasm/crudp"
)

// module wraps a handler for site registration
type module struct {
	handler any
	name    string
}

// siteHandler manages the internal state of the site
type siteHandler struct {
	registeredModules []*module
	DevMode           bool
	cp                *crudp.CrudP
}

var (
	// Singleton instances
	handler = &siteHandler{
		cp: crudp.New(),
	}
)
