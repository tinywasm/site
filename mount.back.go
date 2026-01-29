//go:build !wasm

package site

import (
	"net/http"
	"os"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/client"
	"github.com/tinywasm/fmt"
)

func init() {
	for _, arg := range os.Args {
		if arg == "-dev" {
			handler.DevMode = true
			break
		}
	}
}

// Mount configures the server handled by site.
// It initializes assetmin and registers all routes.
func Mount(mux *http.ServeMux) error {
	// Default Configuration
	// We want to be "Zero configSite" for the user.
	// We can check if we are in dev mode via crudp or environment if needed.
	// For now, let's assume we want to be safe and efficient.

	// Create Javascript handler
	jsHandler := &client.Javascript{
		UseTinyGo:    client.ParseUseTinyGoFlag(),
		WasmFilename: "client.wasm",
	}

	jsHandler.RegisterRoutes(mux, "./public/client.wasm")

	// Create AssetMin instance (private/internal)
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: "./public",

		GetSSRClientInitJS: func() (string, error) { return jsHandler.GetSSRClientInitJS() },

		// We can infer DevMode from crudp or assume production by default for safety in non-dev envs?
		// But usually `site` usage suggests a simple server.
		// Let's use `false` (Production/Efficient) by default for new assetmin,
		// or maybe expose a `site.SetDevMode` that propagates.
		// Given crudp has SetDevMode, we should probably check that.
		DevMode: handler.DevMode,
	})

	// Register AssetMin Routes
	am.RegisterRoutes(mux)

	// build Assets (generate/minify)
	if err := build(am); err != nil {
		fmt.Println("site: build error:", err)
		return err
	}

	// Register CrudP Routes
	handler.cp.RegisterRoutes(mux)

	return nil
}
