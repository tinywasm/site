//go:build wasm

package site

import (
	"github.com/tinywasm/fmt"
)

// No init needed - asset registration is handled by SSR (mount.back.go)

// Mount hydrates the initial module and blocks forever.
func Mount(parentID string) error {
	// 1. Initialize Client (CrudP)
	handler.cp.InitClient()

	// 2. Start the site module management
	if err := Start(parentID); err != nil {
		fmt.Printf("Error starting site: %v\n", err)
		return err
	}

	select {} // Block automatically (WASM apps don't exit)
}

// Render registers the site handlers with the provided mux and prepares assets.
// DEPRECATED: Use Mount(parentID) instead.
func Render() error {
	return Mount("app")
}
