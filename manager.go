package site

import (
	"github.com/tinywasm/dom"
)

var (
	activeModule Module
	cache        []Module
	maxCache     = 3
)

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

// Start initializes the site by hydrating the current module.
func Start(parentID string) error {
	hash := dom.GetHash()
	if hash == "" {
		hash = "#home" // Default route
	}

	name := hash[1:] // Remove #
	m := findModule(name)
	if m == nil {
		return dom.Mount(parentID, nil) // Or a 404 component
	}

	activeModule = m
	return dom.Hydrate(parentID, m)
}

// Navigate switches to a different module.
func Navigate(parentID string, name string) error {
	if activeModule != nil && activeModule.HandlerName() == name {
		return nil
	}

	target := findModule(name)
	if target == nil {
		return nil // Or handle 404
	}

	// 1. Unmount current
	if activeModule != nil {
		dom.Unmount(activeModule)
		addToCache(activeModule)
	}

	// 2. Check cache for target
	if cached := getFromCache(name); cached != nil {
		target = cached
	}

	// 3. Mount new
	activeModule = target
	dom.SetHash("#" + name)
	return dom.Mount(parentID, target)
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

func addToCache(m Module) {
	// Simple LRU: remove oldest if full
	for i, cm := range cache {
		if cm.HandlerName() == m.HandlerName() {
			// Already in cache, move to front
			cache = append(cache[:i], cache[i+1:]...)
			break
		}
	}

	if len(cache) >= maxCache {
		cache = cache[1:]
	}
	cache = append(cache, m)
}

func getFromCache(name string) Module {
	for i, m := range cache {
		if m.HandlerName() == name {
			// Move to front (latest)
			cache = append(cache[:i], cache[i+1:]...)
			cache = append(cache, m)
			return m
		}
	}
	return nil
}
