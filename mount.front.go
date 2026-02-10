//go:build wasm

package site

// No init needed - asset registration is handled by SSR (mount.back.go)

// Mount hydrates the initial module and blocks forever.
func Mount(parentID string) error {
	// 1. Initialize Client (CrudP)
	handler.cp.InitClient()

	// 2. Start the site module management
	if err := Start(parentID); err != nil {
		return err
	}

	select {} // Block automatically (WASM apps don't exit)
}

// Render renders the active component according to the current route.
// It defaults to rendering on the element with ID "app".
// DEPRECATED: Use Mount(parentID) instead.
func Render() error {
	// 1. Initialize Client (CrudP)
	// This connects the responses to the handlers
	handler.cp.InitClient()

	// 2. Start the site module management
	return Start("app")
}
