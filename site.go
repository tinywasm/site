package site

import (
	"github.com/tinywasm/crudp"
)

type assetRegister interface {
	add(handlers ...any) error
}

// siteHandler manages the shared state of the site
type siteHandler struct {
	DevMode           bool
	cp                *crudp.CrudP
	registeredModules []*registeredModule
}

// registeredModule wraps a handler for site registration
type registeredModule struct {
	handler any
	name    string
}

var (
	// Singleton instances
	handler = &siteHandler{
		cp: crudp.New(),
	}
)

// SetUserRoles configures the function to extract user roles from the request context.
func SetUserRoles(fn func(data ...any) []byte) {
	handler.cp.SetUserRoles(fn)
}

// SetAccessCheck configures an external access check function.
// When set, AllowedRoles() interface is NOT required on handlers.
// Must be called before RegisterModules().
//
// Integration with rbac:
//
//	site.SetAccessCheck(func(resource string, action byte, data ...any) bool {
//	    for _, d := range data {
//	        if req, ok := d.(*http.Request); ok {
//	            userID := req.Header.Get("X-User-ID")
//	            ok, _ := rbac.HasPermission(userID, resource, action)
//	            return ok
//	        }
//	    }
//	    return false
//	})
func SetAccessCheck(fn func(resource string, action byte, data ...any) bool) {
	handler.cp.SetAccessCheck(fn)
}

func (h *siteHandler) GetUserData() (name, area string) {
	for _, m := range h.registeredModules {
		if prov, ok := m.handler.(interface {
			GetUserData() (name, area string)
		}); ok {
			n, a := prov.GetUserData()
			if n != "" && a != "" {
				return n, a
			}
		}
	}
	return "Usuario", "Area"
}
