//go:build wasm

package site

import (
	"github.com/tinywasm/dom"
)

// Mount mounts the active component according to the current route.
func Mount(parentID string) error {
	if len(handler.registeredModules) == 0 {
		return nil
	}

	// For now, we just mount the first one as a proof of concept.
	// In a real SPA, this would use a router to pick the component.
	for _, m := range handler.registeredModules {
		if mountable, ok := m.handler.(dom.Mountable); ok {
			return dom.Mount(parentID, mountable)
		}
	}

	return nil
}
