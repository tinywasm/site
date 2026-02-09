//go:build wasm

package site

func registerAssets(handlers ...any) error {
	// No-op in WASM
	return nil
}
