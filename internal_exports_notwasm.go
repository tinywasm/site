//go:build !wasm

package site

import "github.com/tinywasm/assetmin"

// TestSSRBuild exposes the internal ssrBuild function for testing.
// For testing purposes only.
func TestSSRBuild(am *assetmin.AssetMin) error {
	return ssrBuild(am)
}
