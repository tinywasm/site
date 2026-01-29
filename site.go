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

// SetUserRoles configures the function to extract user roles from the request context.
// This is required when using handlers with access control (AllowedRoles).
func SetUserRoles(fn func(data ...any) []byte) {
	handler.cp.SetUserRoles(fn)
}
