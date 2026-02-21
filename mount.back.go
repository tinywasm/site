//go:build !wasm

package site

import (
	"net/http"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/client"
	"github.com/tinywasm/fmt"
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
	if err := applyRBAC(); err != nil {
		return err
	}
	if rbacInitialized && getUserID == nil {
		return fmt.Err("site: SetUserID must be called when using SetDB")
	}
	if !rbacInitialized && !config.DevMode {
		return fmt.Err("site: security not configured — call SetDB or set APP_ENV=development")
	}

	// Create Javascript handler
	jsHandler := client.NewJavascriptFromArgs()

	jsHandler.RegisterRoutes(mux, config.OutputDir+"/client.wasm")

	// Create AssetMin instance
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir:          config.OutputDir,
		GetSSRClientInitJS: func() (string, error) { return jsHandler.GetSSRClientInitJS() },
		DevMode:            config.DevMode,
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
