//go:build wasm

package site

import (
	"github.com/tinywasm/dom"
)

// Mount mounts the active component according to the current route.
// It defaults to mounting on the element with ID "app".
func Mount() error {
	if len(handler.registeredModules) == 0 {
		return nil
	}

	// 1. Mount Active Component
	// In a real SPA, this would use a router to pick the component.
	// For now, we defaults to "app", users can change this in the future if we add config.
	targetID := "app"

	mounted := false
	for _, m := range handler.registeredModules {
		if mountable, ok := m.handler.(dom.Mountable); ok {
			if err := dom.Mount(targetID, mountable); err != nil {
				return err
			}
			mounted = true
			break // Mount only the first one for now
		}
	}

	if !mounted {
		return nil
	}

	// 2. Initialize Client (CrudP)
	// This connects the responses to the handlers
	handler.cp.InitClient()

	return nil
}
