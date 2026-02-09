//go:build wasm

package site

// SetColorPalette is a no-op on WASM as CSS injection happens on the server.
func SetColorPalette(p ColorPalette) {
	// No-op
}
