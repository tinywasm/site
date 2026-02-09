//go:build !wasm

package site

func registerAssets(handlers ...any) error {
	if ssr.assetRegister == nil {
		return nil
	}
	return ssr.assetRegister.add(handlers...)
}
