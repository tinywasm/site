//go:build wasm

package site

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

var (
	activeModule Module
	cache        []Module
)

// Start initializes the site by hydrating the current module.
func Start(parentID string) error {
	hash := dom.GetHash()
	moduleName, params := parseRoute(hash)

	m := findModule(moduleName)
	if m == nil {
		return fmt.Errf("module not found: %s", moduleName)
	}

	// Set params
	if p, ok := m.(Parameterized); ok {
		p.SetParams(params)
	}

	activeModule = m

	if err := dom.Render(parentID, m); err != nil {
		return err
	}

	// Call AfterNavigateTo hook
	if lc, ok := m.(ModuleLifecycle); ok {
		lc.AfterNavigateTo()
	}

	return nil
}

// Navigate switches to a different module based on the hash.
func Navigate(parentID string, hash string) error {
	moduleName, params := parseRoute(hash)

	if activeModule != nil && activeModule.HandlerName() == moduleName {
		// Same module, just update params
		if p, ok := activeModule.(Parameterized); ok {
			p.SetParams(params)
		}

		// Call AfterNavigateTo hook as params changed
		if lc, ok := activeModule.(ModuleLifecycle); ok {
			lc.AfterNavigateTo()
		}
		return nil
	}

	target := findModule(moduleName)
	if target == nil {
		return nil // Or handle 404
	}

	// 1. Check if current module allows navigation away
	if activeModule != nil {
		if lc, ok := activeModule.(ModuleLifecycle); ok {
			if !lc.BeforeNavigateAway() {
				return nil // Cancelled
			}
		}
		// dom.Render handles unmount of previous content automatically
		addToCache(activeModule)
	}

	// 2. Check cache for target
	if cached := getFromCache(moduleName); cached != nil {
		target = cached
	}

	// 3. Set params on new module
	if p, ok := target.(Parameterized); ok {
		p.SetParams(params)
	}

	// 4. Mount new module
	activeModule = target
	dom.SetHash(hash)
	if err := dom.Render(parentID, target); err != nil {
		return err
	}

	// 5. Call AfterNavigateTo hook
	if lc, ok := target.(ModuleLifecycle); ok {
		lc.AfterNavigateTo()
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

	if len(cache) >= config.CacheSize {
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
