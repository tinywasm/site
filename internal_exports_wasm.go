//go:build wasm

package site

// TestResetWasm resets the active module and cache for testing.
// For testing purposes only.
func TestResetWasm() {
	activeModule = nil
	cache = nil
}
