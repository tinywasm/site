//go:build !wasm

package site

import (
	"net/http"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/fmt"
)

// Config defines the configuration for the site mount
type Config struct {
	PublicDir   string
	DevMode     bool
	AssetsCache bool // Forces assetmin to cache or not. If nil/false/true depends on logic.
}

// Mount configures the server handled by site.
// It initializes assetmin and registers all routes.
func Mount(mux *http.ServeMux) error {
	// Default Configuration
	// We want to be "Zero Config" for the user.
	// We can check if we are in dev mode via crudp or environment if needed.
	// For now, let's assume we want to be safe and efficient.

	// Create AssetMin instance (private/internal)
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: "./public",
		// We can infer DevMode from crudp or assume production by default for safety in non-dev envs?
		// But usually `site` usage suggests a simple server.
		// Let's use `false` (Production/Efficient) by default for new assetmin,
		// or maybe expose a `site.SetDevMode` that propagates.
		// Given crudp has SetDevMode, we should probably check that.
		DevMode: cp.IsDevMode(), // Make sure crudp exposes IsDevMode or we just track it.
	})

	// Register AssetMin Routes
	am.RegisterRoutes(mux)

	// Build Assets (generate/minify)
	if err := Build(am); err != nil {
		fmt.Println("site: build error:", err)
		return err
	}

	// Register CrudP Routes
	cp.RegisterRoutes(mux)

	return nil
}
