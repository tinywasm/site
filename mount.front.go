//go:build wasm

package site

// No init needed - asset registration is handled by SSR (mount.back.go)

// Render renders the active component according to the current route.
// It defaults to rendering on the element with ID "app".
func Render() error {
	// 1. Initialize Client (CrudP)
	// This connects the responses to the handlers
	handler.cp.InitClient()

	// 2. Start the site module management
	return Start("app")
}
