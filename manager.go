package site

import (
	"strings"
)

// parseRoute extracts module name and params from hash
func parseRoute(hash string) (module string, params []string) {
	if hash == "" || hash == "#" {
		return config.DefaultRoute, nil // Default route
	}

	cleanHash := strings.TrimPrefix(hash, "#")
	// Remove leading slash if any to support #/users style
	cleanHash = strings.TrimPrefix(cleanHash, "/")

	if cleanHash == "" {
		return config.DefaultRoute, nil
	}

	parts := strings.Split(cleanHash, "/")
	if len(parts) == 0 {
		return config.DefaultRoute, nil
	}

	return parts[0], parts[1:]
}

// registerModule adds a module to the site registry.
func registerModule(m Module) {
	name := m.HandlerName()
	exists := false
	for _, rm := range handler.registeredModules {
		if rm.name == name {
			exists = true
			break
		}
	}
	if !exists {
		handler.registeredModules = append(handler.registeredModules, &registeredModule{
			handler: m,
			name:    name,
		})
	}
}

func findModule(name string) Module {
	for _, rm := range handler.registeredModules {
		if rm.name == name {
			if m, ok := rm.handler.(Module); ok {
				return m
			}
		}
	}
	return nil
}
