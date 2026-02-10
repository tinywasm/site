//go:build !wasm

package site

import (
	"net/http"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/client"
)

type backendRegister struct {
	handlers []any
}

func (r *backendRegister) add(handlers ...any) error {
	r.handlers = append(r.handlers, handlers...)
	return nil
}

func init() {
	ssr.assetRegister = &backendRegister{}
}

// Mount registers the site handlers with the provided mux and prepares assets.
func Mount(mux *http.ServeMux) error {
	// Default Configuration
	// We want to be "Zero configSite" for the user.
	// We can check if we are in dev mode via crudp or environment if needed.
	// For now, let's assume we want to be safe and efficient.

	// Create Javascript handler
	jsHandler := client.NewJavascriptFromArgs()

	jsHandler.RegisterRoutes(mux, config.OutputDir+"/client.wasm")

	// Create AssetMin instance (private/internal)
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: config.OutputDir,

		GetSSRClientInitJS: func() (string, error) { return jsHandler.GetSSRClientInitJS() },

		// We can infer DevMode from crudp or assume production by default for safety in non-dev envs?
		// But usually `site` usage suggests a simple server.
		// Let's use `false` (Production/Efficient) by default for new assetmin,
		// or maybe expose a `site.SetDevMode` that propagates.
		// Given crudp has SetDevMode, we should probably check that.
		DevMode: config.DevMode,
	})

	// Register assets from modules
	// Handled by ssrBuild(am) which processes handler.registeredModules

	// ssrBuild Assets (generate/minify) - MUST happen BEFORE RegisterRoutes
	// to ensure the sprite is complete before accepting requests
	if err := ssrBuild(am); err != nil {
		return err
	}

	// Register AssetMin Routes AFTER ssrBuild to ensure sprite is complete
	am.RegisterRoutes(mux)

	// Register CrudP Routes
	handler.cp.RegisterRoutes(mux)

	return nil
}

// Render registers the site handlers with the provided mux and prepares assets.
// DEPRECATED: Use Mount(mux) instead.
func Render(mux *http.ServeMux) error {
	return Mount(mux)
}
